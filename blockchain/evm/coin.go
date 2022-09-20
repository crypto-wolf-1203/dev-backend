package evm

type CoinInfo struct {
	Name string
	ChainName string
	Address string
}

var coniMap map[string]CoinInfo = map[string]CoinInfo{
	"BNB": {
		Name: "BNB",
		ChainName: "Binance Smart Chain",
		Address: "0x10ED43C718714eb63d5aA57B78B54704E256024E",
	},
	"BUSD": {
		Name: "BUSD",
		ChainName: "Binance Smart Chain",
		Address: "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56",
	},
	"USDT": {
		Name: "USDT",
		ChainName: "Binance Smart Chain",
		Address: "0x55d398326f99059fF775485246999027B3197955",
	},
	"USDC": {
		Name: "USDC",
		ChainName: "Binance Smart Chain",
		Address: "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d",
	},
}
