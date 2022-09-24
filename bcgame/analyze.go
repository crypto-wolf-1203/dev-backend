package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"strconv"
)

func main() {
	file, err := os.Open("1.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)

	inputData := make([]int, 0)

	var min int = 100000000000
	var max int = 0

	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			f := val
			inputData = append(inputData, f)
			if min > f {
				min = f
			}
			if max < f {
				max = f
			}
		}
	}

	// reversing data from file
	v := make([]int, len(inputData))
	for idx, val := range inputData {
		v[len(inputData) - 1 - idx] = val
	}

	// analyzing data
	var hist map[int]int = map[int]int{}

	for _, score := range v {
		h, ok := hist[score]
		if ok {
			hist[score] = h + 1
		} else {
			hist[score] = 1
		}
	}

	// fmt.Println("")
	// fmt.Println("*****************************")
	// fmt.Println("histogram")

	// for k1, v1 := range hist {
	// 	fmt.Println(k1, v1)
	// }
	
	// fmt.Println("")
	// fmt.Println("*****************************")
	// fmt.Println("z values")

	for z := min; z <= 3000; z += 1.0 {
		zxmax_integral := 0
		for k := z; k <= max; k += 1.0 {
			vv, ok := hist[k]
			if ok {
				zxmax_integral += vv
			}
		}

		xminz_integral := 0
		for k := min; k < z; k += 1.0 {
			vv, ok := hist[k]
			if ok {
				xminz_integral += vv
			}
		}

		f := (float64(z) * float64(zxmax_integral) / 100.0 - float64(xminz_integral)) / float64(len(v))
		fmt.Printf("%v,%v,%v,%v\n", z, f, float64(zxmax_integral) / float64(len(v)), float64(xminz_integral) / float64(len(v)))
	}
}