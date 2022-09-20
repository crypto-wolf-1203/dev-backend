package osdep

import (
	"os"
)

func Getenv(key string, alt string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return alt
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
