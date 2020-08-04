package service

import (
	"math/rand"
)

func (sc *ServiceConfig) randomURL() string {
	return sc.ChainRemoteURL[rand.Intn(len(sc.ChainRemoteURL))]
}
