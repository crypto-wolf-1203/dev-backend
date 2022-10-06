package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"pongpongi.com/blockchain/evm"
)

func ProcGet(subURL string, query map[string]interface{}, w http.ResponseWriter) error {
	var err error
	err = errors.New("unexpected result")
	switch subURL {
	case "/read":
		coinArray, existing := query["price"]
		if existing {
			switch reflect.TypeOf(coinArray).Kind() {
			case reflect.Slice:
				// fmt.Println("Slice received")
			case reflect.Array:
				// fmt.Println("Array received")
			}

			coinLabel := coinArray.([]string)[0]
			reportCoinPrice(coinLabel, w)
		} else {
			return err
		}
	case "/write":
	default:
		return err
	}
	return nil
}

// https://pongpongi.com/api/read?price=BNB
func reportCoinPrice(coin string, w http.ResponseWriter) error {
	coinInfo, err1 := evm.GetCoinInfo(coin)
	usdtInfo, errUSDT := evm.GetCoinInfo("USDT");

	if err1 != nil {
		return err1
	} else if errUSDT != nil {
		return errUSDT
	}

	coinPrice, err := evm.GetPairRatio(coinInfo.Address, usdtInfo.Address)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%v", coinPrice)
	return nil
}
