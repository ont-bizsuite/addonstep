package service

import (
	ontsdk "github.com/ontio/ontology-go-sdk"
)

type ServiceConfig struct {
	RecvAddress      string
	ONGFloatPriceURL string
	NetType          string
	GasPrice         uint64
	GasLimit         uint64
	ChainRemoteURL   []string
	Account          *ontsdk.Account // addon fee transfer to address
	// ons related
	ONSContractAddress string
	ONSUpLevelAccount  *ontsdk.Account // normally this is the same as collect money address, when register ons, one must provide the the upper level permission
}

type ServiceConfigOption func(*ServiceConfig)

func RecvAddressOption(rcvAddr string) ServiceConfigOption {
	return func(pc *ServiceConfig) {
		pc.RecvAddress = rcvAddr
	}
}

func ONGFloadPriceOption(furl string) ServiceConfigOption {
	return func(pc *ServiceConfig) {
		pc.ONGFloatPriceURL = furl
	}
}

func NetTypeOption(nt string) ServiceConfigOption {
	return func(pc *ServiceConfig) {
		pc.NetType = nt
	}
}

func GasPriceOption(gp uint64) ServiceConfigOption {
	return func(pc *ServiceConfig) {
		pc.GasPrice = gp
	}
}

func GasLimitOption(gl uint64) ServiceConfigOption {
	return func(pc *ServiceConfig) {
		pc.GasLimit = gl
	}
}

func ChainRemoteURLOption(rs []string) ServiceConfigOption {
	return func(pc *ServiceConfig) {
		pc.ChainRemoteURL = rs
	}
}

func WalletOption(w *ontsdk.Account) ServiceConfigOption {
	return func(pc *ServiceConfig) {
		pc.Account = w
	}
}

func NewConfig(opts ...ServiceConfigOption) *ServiceConfig {
	ret := &ServiceConfig{}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
