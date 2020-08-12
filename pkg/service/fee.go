package service

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ont-bizsuite/addonstep/pkg/path"
	ontsdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/utils"
)

var (
	recentSuccessONGPrice float64
)

// PayParam
type PayParam struct {
}

var (
	PaySample = &PayParam{}
)

func (pc *ServiceConfig) TransferMoney(c *gin.Context) {
	tmpl, err := pc.payme()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to generate config",
		})
	}

	c.Data(200, "text/plain", []byte(tmpl))
}

var (
	client = http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

type FloatONGResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result Result `json:"result"`
}
type USD struct {
	Price      string `json:"price"`
	Percentage string `json:"percentage"`
}
type Prices struct {
	USD USD `json:"USD"`
}
type Result struct {
	Token  string `json:"token"`
	Rank   int    `json:"rank"`
	Prices Prices `json:"prices"`
}

func decodeFloatOng(r io.Reader, param *FloatONGResp) (float64, error) {
	err := json.NewDecoder(r).Decode(param)
	if err != nil {
		return 0.0, err
	}

	if param.Code != 0 {
		return 0.0, errors.New("response not ok, expect 0")
	}

	return strconv.ParseFloat(param.Result.Prices.USD.Price, 64)
}

func (pc *ServiceConfig) GetFloatONGPrice() (float64, error) {
	req, err := http.NewRequest(http.MethodGet, pc.ONGFloatPriceURL, nil)
	if err != nil {
		return 0.0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()
	var param FloatONGResp

	return decodeFloatOng(resp.Body, &param)
}

func (pc *ServiceConfig) dollarToONG() int64 {
	rtPrice, err := pc.GetFloatONGPrice()
	var fail bool
	if err != nil {
		log.Printf("can not get real time ong price %v", err)
		// price of 20200701
		rtPrice = recentSuccessONGPrice
		fail = true
	}
	if rtPrice == 0 {
		rtPrice = 0.0000000001
	}
	if !fail {
		recentSuccessONGPrice = rtPrice
	}

	// 3 month fee is 50 * 3
	return int64(150 / rtPrice)
}

func (pc *ServiceConfig) paydata(ong int64) (string, error) {
	var ret TxPrepareData
	ret.Action = "signTransaction"
	ic := &InvokeConfig{}
	ic.GasPrice = pc.GasPrice
	ic.GasLimit = pc.GasLimit
	ic.Payer = "%address"
	ic.ContractHash = "0200000000000000000000000000000000000000"

	ic.Functions = []Functions{
		Functions{
			Operation: "transfer",
			Args: []Args{
				Args{
					Name:  "from",
					Value: FromString("Address:%address"),
				},
				Args{
					Name:  "to",
					Value: FromString(fmt.Sprintf("Address:%s", pc.Account.Address.ToBase58())),
				},
				Args{
					Name:  "amount",
					Value: FromInt(ong * 1e9),
				},
			},
		},
	}

	p := &Params{
		InvokeConfig: ic,
	}
	ret.Params = p

	b, _ := json.Marshal(ret)
	return string(b), nil
}

func (pc *ServiceConfig) id() string {
	return "String:did:ont:" + pc.Account.Address.ToBase58()
}

func (pc *ServiceConfig) signData(data []byte) string {
	sign, err := pc.Account.Sign(data)
	if err != nil {
		return ""
	}

	// TODO: explaim 01 means, just add this to make mobile app happy
	return "01" + hex.EncodeToString(sign)
}

func (pc *ServiceConfig) payme() (string, error) {
	var ret TxPrepare
	ong := pc.dollarToONG()
	ret.Desc = Desc{
		Price:  fmt.Sprintf("%dong", ong),
		Detail: "transfer ong to addon store",
		Type:   "transfer ong to addon store",
	}
	ret.Signer = ""
	ret.Exp = time.Now().Add(time.Minute * 10).Unix()
	ret.ID = "111"
	ret.Callback = path.PayCallbackPath
	// ret.Data = "todo"
	ret.Chain = pc.NetType
	ret.Ver = "v2.0.0"
	data, err := pc.paydata(ong)
	if err != nil {
		return "", err
	}
	ret.Data = data

	ret.Requester = pc.id()
	ret.Signature = pc.signData([]byte(ret.Data))

	b, _ := json.Marshal(ret)

	return string(b), nil
}

type SendTxParam struct {
	SingedTx string      `json:"signed_tx"`
	Prev     interface{} `json:"prev"`
}

func (pc *ServiceConfig) FeeBack(c *gin.Context) {
	var param SendTxParam
	err := c.BindJSON(&param)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("raw tx: %s", param.SingedTx)

	sdk := ontsdk.NewOntologySdk()
	node := pc.ChainRemoteURL[rand.Intn(len(pc.ChainRemoteURL))]
	sdk.NewRpcClient().SetAddress(node)
	log.Printf("send to chain: %s", node)

	tx, err := utils.TransactionFromHexString(param.SingedTx)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "invalid tx",
		})
		return
	}
	mtx, err := tx.IntoMutable()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "invalid tx, into mutable fail",
		})
		return
	}
	hash, err := sdk.SendTransaction(mtx)
	log.Printf("fee tx hash: %s", hash.ToHexString())

	if err != nil {
		log.Printf("tx fail: %v", err)
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("tx fai;: %v", err),
		})
		return
	}

	if err := pollEvent(hash.ToHexString(), sdk); err != nil {
		log.Printf("poll tx event fail: %v", err)
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("tx fai;: %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "addon fee ok",
	})
}

func pollEvent(hashHex string, sdk *ontsdk.OntologySdk) error {
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)

		event, err := sdk.GetSmartContractEvent(hashHex)
		if err != nil {
			return err
		}

		if event == nil {
			continue
		}

		if event.State != 1 {
			return errors.New("event state not 1")
		}
		// get the state is 1
		return nil
	}

	return errors.New("try 5 times and still get no states")
}
