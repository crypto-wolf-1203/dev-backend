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

	command := args[0]
	if command == "fixed" {
		gap, _ := strconv.Atoi(args[1])
		keep, _ := strconv.Atoi(args[2])
		rate, _ := strconv.Atoi(args[3])

		rate -= 100
		if rate < 0 {
			panic("error setting lower than 100")
		}

		fmt.Println(gap, keep, rate)

		// scanning through
		method1(rate, gap, keep, v)
	} else if command == "double" {
		gap, _ := strconv.Atoi(args[1])
		keep, _ := strconv.Atoi(args[2])
		doubleMethod(gap, keep, v)
	} else {
		fmt.Println("undefined method " + command)
	}
}

// go run sim.go fixed <gap> <keep> <rate+100>
// go run sim.go fixed 0 40 10100
func method1(rate, gap, keep int, v []int) {
	skip := -1
	sum := 0
	deductCounter := 0
	loseDepth := 0
	maxLoseDepth := 0
	rewindTimes := 0
	profitTimes := 0
	missedTimes := 0
	minSum := 10000000000
	defBetOnFirstDeduct := true
	betOnDeduct := defBetOnFirstDeduct
	deductTimes := 0
	betRet := false

	for _, val := range v {
		fmt.Printf("+++ rate=%v, bet=%v +++\n", val+100, betRet)

		if betRet {
			if val >= rate {
				sum += rate
				profitTimes++
				if maxLoseDepth < loseDepth {
					maxLoseDepth = loseDepth
				}
				loseDepth = 0
				fmt.Println("*************** profitting", sum)
			} else {
				sum -= 100
				deductTimes++
				if minSum > sum {
					minSum = sum
				}
				loseDepth++
				fmt.Println("lost", sum)
			}
		}

		betOnDeduct = false

		if skip < 0 && val >= rate {
			skip = 0
			rewindTimes++
			if gap == 0 {
				betOnDeduct = defBetOnFirstDeduct
			}
			fmt.Println("initialized")
		} else if skip < 0 {
			fmt.Println("skipped for initialization")
		} else if skip < gap {
			skip++
			fmt.Println("skipped", skip)
			if skip == gap {
				betOnDeduct = defBetOnFirstDeduct
			}
		} else if deductCounter < keep {
			if betRet && val >= rate {
				skip = 0
				rewindTimes++
				if gap == 0 {
					betOnDeduct = defBetOnFirstDeduct
				}
				deductCounter = 0
				fmt.Println("rewinded by profitting")
			} else {
				// if betRet {
				// 	fmt.Println("skipped deducting...")
				// } else {
				betOnDeduct = true
				deductCounter++
				// }
			}
		} else {
			if val >= rate {
				skip = 0
				rewindTimes++

				if !betRet {
					missedTimes++
				}

				if gap == 0 {
					betOnDeduct = defBetOnFirstDeduct
				}
				deductCounter = 0
				fmt.Println("rewinded")
			} else {
				fmt.Println("ignoring")
			}
		}

		betRet = betOnDeduct
	}

	fmt.Println("")
	fmt.Println(sum/100, "Total sum")
	fmt.Println(minSum/100, "Lowest sum")
	fmt.Println(maxLoseDepth, "lose depth maximum")
	fmt.Println(rewindTimes, "deduct period count")
	fmt.Println(deductTimes, "deduct count")
	fmt.Println(profitTimes, "profit period count")
	fmt.Println(missedTimes, "missed period count")
}

// go run sim.go double <gap> <max-keep>
// go run sim.go double

func doubleMethod(gap, maxKeep int, v []int) {
	skip := -1
	betThisTime := false
	amount := 0
	sum := 0
	keep := 0
	mode := 2 // 0 :skipping, 1: deducting, 2: ignoring
	failCount := 0
	profitCount := 0
	minSum := 10000000
	deductCount := 0
	minRound := 100000
	minRoundIdx := -1
	maxRound := 0
	maxRoundIdx := -1

	roundCount := 0
	avgRound := 0

	intervalIdx := make([]int, 0)
	intervalLength := make([]int, 0)

	dominanceLimit := 64

	for idx, val := range v {
		fmt.Printf("[%v, %v] ", idx, val)

		if betThisTime && val >= 100 {
			if profitCount+deductCount < dominanceLimit {
				sum += amount * 2
				fmt.Print(" *** profitting *** ")
			}
			profitCount++
			deductCount--
		}

		nextMode := 0

		switch mode {
		case 0:
			if skip+1 == gap {
				nextMode = 1
			} else {
				nextMode = mode
				skip++
				if val >= 100 {
					skip = 0
				}
			}
		case 1:
			if keep < maxKeep {
				nextMode = mode
				if val >= 100 {
					nextMode = 0
				}
				keep++
			} else {
				nextMode = 2
				if betThisTime && val < 100 {
					if profitCount+deductCount < dominanceLimit {
						failCount++
					}
					if profitCount+deductCount < minRound {
						minRound = profitCount + deductCount
						minRoundIdx = idx
					}

					if profitCount+deductCount > maxRound {
						maxRound = profitCount + deductCount
						maxRoundIdx = idx
					}

					intervalIdx = append(intervalIdx, idx)
					intervalLength = append(intervalLength, profitCount+deductCount)

					roundCount++
					avgRound += deductCount + profitCount
				}

				if val >= 100 {
					nextMode = 0
				}
			}
		case 2:
			if val >= 100 {
				nextMode = 0
			} else {
				nextMode = mode
			}
		}

		if mode != nextMode {
			keep = 0
			skip = 0
			amount = 0
			betThisTime = false
		}

		switch nextMode {
		case 0:
			if gap == 0 {
				betThisTime = true
				nextMode = 1
				amount = 1
				keep = 1
			}
		case 1:
			betThisTime = true
			if amount == 0 {
				amount = 1
			} else {
				amount *= 2
			}
		case 2:
			fmt.Printf(" --- failing --- deduct %d, profit %d, ", deductCount, profitCount)
			deductCount = 0
			profitCount = 0
		}

		if betThisTime {
			if profitCount+deductCount < dominanceLimit {
				sum -= amount
			}
			deductCount++
		}

		if sum < minSum {
			minSum = sum
		}

		avg := 0
		if roundCount > 0 {
			avg = avgRound / roundCount
		}
		fmt.Printf("mode %d, next mode %d, bet %v, amount %v, sum %d, worst %d, fail %d, profit %d, deduct %d, min round %d (index %d), max round %d (index %d), avg round %d\n", mode, nextMode, betThisTime, amount, sum, minSum, failCount, profitCount, deductCount, minRound, minRoundIdx, maxRound, maxRoundIdx, avg)

		mode = nextMode
	}

	fmt.Println("interval list:")
	for idx, val := range intervalIdx {
		fmt.Printf("%d, %d\n", val, intervalLength[idx])
	}
}
