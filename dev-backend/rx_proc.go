package main

import (
	"fmt"
	"errors"
	"reflect"
)

func ProcGet(subURL string, query map[string]interface{}) error {
	var err error
	err = errors.New("unexpected result")
	switch subURL {
	case "/read":
		coinArray, existing := query["price"]
		if !existing {
			return err
		}
		switch reflect.TypeOf(coinArray).Kind() {
		case reflect.Slice:
			fmt.Println("Slice received")
		case reflect.Array:
			fmt.Println("Array received")
		}

		coinLabel := coinArray.([]string)[0]
		fmt.Println(coinLabel)
	case "/write":
		break
	default:
		return err
	}
	return nil
}
