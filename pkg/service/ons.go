package service

import (
	// "encoding/json"
	// "errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	ontsdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	// "github.com/rongyi/triones-node/pkg/wallet"
)

var (
	RegisterONSSample = &RegisterONSParam{
		Domain: "",
	}
)

// RegisterONSParam is the param for create ons
// Domain: the domain to be registered
// OntID: owner of this domain
type RegisterONSParam struct {
	Domain string      `json:"Domain"`
	OntID  string      `json:"ontid,omitempty"`
	Prev   interface{} `json:"prev,omitempty"`
}

func (sc *ServiceConfig) RegisterONS(c *gin.Context) {
	var param RegisterONSParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	contracAddr, _ := common.AddressFromHexString(sc.ONSContractAddress)

	params := []interface{}{
		"registerDomain",
		[]interface{}{
			fmt.Sprintf("%s.node.ont", param.Domain),
			param.OntID,
			1,
			time.Now().Add(86400000 * time.Second).Unix(),
		},
	}

	sdk := ontsdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(sc.randomURL())

	txHash, err := sdk.NeoVM.InvokeNeoVMContract(sc.GasPrice, sc.GasLimit, sc.ONSUpLevelAccount, sc.ONSUpLevelAccount, contracAddr, params)
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("fail to create domain: %v", err),
		})
		return
	}

	log.Printf("txhash: %s", txHash.ToHexString())

	if err := pollEvent(txHash.ToHexString(), sdk); err != nil {
		log.Printf("fail to get event: %v", err)

		c.JSON(500, gin.H{
			"message": fmt.Sprintf("fail to get event: %v", err),
		})

		return
	}

	c.JSON(200, param)
}
