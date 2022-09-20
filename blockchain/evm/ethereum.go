package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rpc"
	"pongpongi.com/osdep"
	"strings"
	"bufio"
	"os"
)

const erc20SimpleABI = `[{"inputs": [{"internalType": "address","name": "","type": "address"}],"name": "userInfos","outputs": [{"internalType": "uint256","name": "","type": "uint256"},{"internalType": "uint256","name": "","type": "uint256"}],"stateMutability": "view","type": "function"}]`

var bscRPC *rpc.Client
var erc20ABI abi.ABI

func InitEVM() {
	bscRPC = loadRPC("BSC_RPC", "https://bsc-dataseed.binance.org/")
	erc20ABI = loadABIFromFile("./erc20.json")
}

func loadRPC(envName string, defURL string) *rpc.Client {
	rpcInst, err := rpc.DialHTTP(osdep.Getenv(envName, defURL))
	osdep.Check(err)

	return rpcInst
}

func loadABIFromString(abiText string) abi.ABI {
	retABI, err := abi.JSON(strings.NewReader(abiText))
	osdep.Check(err)

	return retABI
}

func loadABIFromFile(abiPath string) abi.ABI {
	f, err := os.Open(abiPath)
	defer f.Close()

	osdep.Check(err)

	reader := bufio.NewReader(f)

	retABI, err2 := abi.JSON(reader)
	osdep.Check(err2)

	return retABI
}
