package evm

import (
	"bufio"
	"os"
	"strings"
  // "fmt"
  "math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rpc"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/common/hexutil"
	"pongpongi.com/osdep"
)

const ADDRESS_ZERO = "0x0000000000000000000000000000000000000000"

const erc20SimpleABI = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transfer",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`

var bscRPC *rpc.Client
var erc20ABI, routerABI, factoryABI abi.ABI

func InitEVM() {
	bscRPC = loadRPC("BSC_RPC", "https://bsc-dataseed.binance.org/")
	erc20ABI = loadABIFromFile("./erc20.json")
  routerABI = loadABIFromFile("./IPancakeRouter02.json")
  factoryABI = loadABIFromFile("./IPancakeFactory.json")
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

func queryContract(contractAddress string, txData []byte, processFunc func(string, interface{}) error, param interface{}) error {
  var resultStr string
  err := bscRPC.Call(&resultStr, "eth_call", map[string]interface{}{
          "from": ADDRESS_ZERO,
          "to":   contractAddress,
          "data": hexutil.Bytes(txData),
        }, "latest")

  if err == nil {
    return processFunc(resultStr, param)
  } else {
    return err
  }
}

func getPairProcess(queryResult string, param interface{}) error {
  result, err := factoryABI.Unpack("getPair", hexutil.MustDecode(queryResult))
  if err != nil {
    return err
  }

  pairAddr := result[0].(common.Address).String()
  *param.(*string) = pairAddr
  return nil
}

func getBalanceAmount(queryResult string, param interface{}) error {
  result, err := erc20ABI.Unpack("balanceOf", hexutil.MustDecode(queryResult))
  if err != nil {
    return err
  }

  *param.(*big.Int) = *result[0].(*big.Int)

  return nil
}

func getDecimalAmount(queryResult string, param interface{}) error {
  result, err := erc20ABI.Unpack("decimals", hexutil.MustDecode(queryResult))
  if err != nil {
    return err
  }

  *param.(*uint8) = result[0].(uint8)

  return nil
}

func GetPairRatio(token1Address, token2Address string) (float64, error) {
  var pairAddress string

  pairData, pairErr := factoryABI.Pack("getPair", common.HexToAddress(token1Address), common.HexToAddress(token2Address))

  if pairErr != nil {
    return 0, pairErr
  }

  qerr := queryContract("0xca143ce32fe78f1f7019d7d551a6402fc5350c73", pairData, getPairProcess, &pairAddress) // pancakeswap factory on bsc mainnet
  if qerr != nil {
    return 0, qerr
  }
  // fmt.Println("********************** " + pairAddress)

  balanceData, balanceErr := erc20ABI.Pack("balanceOf", common.HexToAddress(pairAddress))

  if balanceErr != nil {
    return 0, balanceErr
  }

  decimalData, decimalErr := erc20ABI.Pack("decimals")

  if decimalErr != nil {
    return 0, decimalErr
  }

  var token1Balance, token2Balance big.Int
  var token1Decimal, token2Decimal uint8

  qerr = queryContract(token1Address, balanceData, getBalanceAmount, &token1Balance)
  // fmt.Println("******** " + token1Address, token1Balance.String())
  qerr = queryContract(token1Address, decimalData, getDecimalAmount, &token1Decimal)
  // fmt.Println("******** " + token1Address, token1Decimal)

  qerr = queryContract(token2Address, balanceData, getBalanceAmount, &token2Balance)
  // fmt.Println("******** " + token2Address, token2Balance.String())
  qerr = queryContract(token2Address, decimalData, getDecimalAmount, &token2Decimal)
  // fmt.Println("******** " + token2Address, token2Decimal)

  res := new(big.Float).SetInt(&token2Balance)

  res.Mul(res, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token1Decimal)), nil)))
  res.Quo(res, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token2Decimal)), nil)))
  res.Quo(res, new(big.Float).SetInt(&token1Balance))

  retVal, _ := res.Float64()
  // fmt.Println("======= ", retVal)

  return retVal, nil
}
