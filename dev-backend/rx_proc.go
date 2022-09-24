package main

import (
	"fmt"
	"errors"
	"reflect"
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
				fmt.Println("Slice received")
			case reflect.Array:
				fmt.Println("Array received")
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

func reportCoinPrice(coin string, w http.ResponseWriter) error {
	
}
