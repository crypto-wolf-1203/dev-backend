package evm

import (
	"errors"
)

type ChainInfo struct {
	Name string
	NativeCurrency string
	Decimal int
	ChainId int
	RpcURL string
	BlockExplorer string
}

var chainMap map[string]ChainInfo = map[string]ChainInfo {
	"bsc": {
		Name: "Binance Smart Chain",
		NativeCurrency: "BNB",
		Decimal: 18,
		ChainId: 56,
		RpcURL: "https://bsc-dataseed.binance.org/",
		BlockExplorer: "https://bscscan.com/",
	},
	"bsctestnet": {
		Name: "Binance Smart Chain Testnet",
		NativeCurrency: "TBNB",
		Decimal: 18,
		ChainId: 97,
		RpcURL: "https://data-seed-prebsc-2-s1.binance.org:8545/",
		BlockExplorer: "https://testnet.bscscan.com/",
	},
}

func GetChainInfo(name string) (*ChainInfo, error) {
	value, ok := chainMap[name]
	if ok {
		return &value, nil
	} else {
		return nil, errors.New(name + ": undefined")
	}
}
