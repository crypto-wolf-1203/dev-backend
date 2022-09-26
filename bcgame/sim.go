package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	gap, _ := strconv.Atoi(args[0])
	keep, _ := strconv.Atoi(args[1])
	rate, _ := strconv.Atoi(args[2])

	fmt.Println(gap, keep, rate)

	rate -= 100
	if rate < 0 {
		panic("error setting lower than 100")
	}

	file, err := os.Open("1.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)

	inputData := make([]int, 0)

	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		} else {
			inputData = append(inputData, val)
		}
	}

	// reversing data from file
	v := make([]int, len(inputData))
	for idx, val := range inputData {
		v[len(inputData)-1-idx] = val
	}

	// scanning through
	skip := -1
	sum := 0
	deductCounter := 0
	loseDepth := 0
	maxLoseDepth := 0
	shotTimes := 0
	profitTimes := 0
	minSum := 10000000000
	for idx0, val := range v {
		idx := idx0 + 1
		if skip < 0 && val >= rate {
			skip = 0
			fmt.Println(idx, val, "initialized", val, ">=", rate)
			continue
		}

		if val < rate {
			if skip < 0 {
				fmt.Println(idx, val, "skipping for initialization")
				continue
			}
			if skip < gap {
				skip++
				fmt.Println(idx, val, "skipping", skip)
				continue
			}
			if deductCounter < keep {
				if deductCounter == 0 {
					shotTimes++
				}

				fmt.Println(idx, val, "deducting for", val, "<", rate)
				sum -= 100
				if minSum > sum {
					minSum = sum
				}
				loseDepth++
			} else {
				fmt.Println(idx, val, "ignoring for", deductCounter, ">=", keep)
			}
			deductCounter++
		} else if skip >= gap {
			skip = 0
			if deductCounter < keep {
				sum += rate

				if maxLoseDepth < loseDepth {
					maxLoseDepth = loseDepth
				}
				loseDepth = 0
				fmt.Println(idx, val, "********************** profitting for", val, ">=", rate)
				profitTimes++
			} else {
				fmt.Println(idx, val, "missing for", val, ">=", rate)
			}

			deductCounter = 0
		} else {
			fmt.Println(idx, val, "missing for", val, ">=", rate, "but too early")
			skip = 0
		}
	}

	fmt.Println("")
	fmt.Println(sum/100, "Total sum")
	fmt.Println(minSum/100, "Lowest sum")
	fmt.Println(maxLoseDepth, "lose depth maximum")
	fmt.Println(shotTimes, "deduct period count")
	fmt.Println(profitTimes, "profit period count")
}

// go run sim.go 20 8 2100 // not good, proved being bad
// go run sim.go 50 40 10100 // good, proving...
