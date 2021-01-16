package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

func rollDie() int {
	return rand.Intn(6) + 1
}

func runSim(goal string, numTrials int, txt string) {
	total := 0
	for i := 0; i < numTrials; i++ {
		var result string
		for j := 0; j < len(goal); j++ {
			result += strconv.Itoa(rollDie())
			if result == goal {
				total++
			}
		}
	}
	fmt.Printf("Actual probability of %s = %f\n", txt, 1.0/math.Pow(6.0, float64(len(goal))))
	estProbability := float64(total) / float64(numTrials)
	fmt.Printf("Estimated probability of %s = %f\n", txt, estProbability)
}

func makeRange(min, max int) []int {
	r := make([]int, max-min)
	for i := range r {
		r[i] = min + i
	}
	return r
}

func sameDate(numPeople int, numSame int) bool {
	possibleDates := makeRange(0, 365)
	/*
		var possibleDates []int
		r1 := makeRange(0, 57)
		possibleDates = append(possibleDates, r1...)
		possibleDates = append(possibleDates, r1...)
		possibleDates = append(possibleDates, r1...)
		possibleDates = append(possibleDates, r1...)
		possibleDates = append(possibleDates, 58)
		r2 := makeRange(59, 366)
		possibleDates = append(possibleDates, r2...)
		possibleDates = append(possibleDates, r2...)
		possibleDates = append(possibleDates, r2...)
		possibleDates = append(possibleDates, r2...)
		r3 := makeRange(180, 270)
		possibleDates = append(possibleDates, r3...)
		possibleDates = append(possibleDates, r3...)
		possibleDates = append(possibleDates, r3...)
		possibleDates = append(possibleDates, r3...)
	*/

	birthdays := [366]int{}
	for p := 0; p < numPeople; p++ {
		birthDate := possibleDates[rand.Intn(len(possibleDates))]
		birthdays[birthDate]++
	}
	max := 0
	for _, b := range birthdays {
		if max < b {
			max = b
		}
	}
	return max >= numSame
}

func birthdayProb(numPeople int, numSame int, numTrials int) float64 {
	numHits := 0
	for t := 0; t < numTrials; t++ {
		if sameDate(numPeople, numSame) {
			numHits++
		}
	}
	return float64(numHits) / float64(numTrials)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	runSim("11111", 1000000, "11111")

	for _, numPeople := range [...]int{10, 20, 40, 100} {
		fmt.Println("For", numPeople, "est. prob. of a shared birthday is", birthdayProb(numPeople, 2, 100000))

		factorial366 := new(big.Int).MulRange(1, 366)
		factorialN := new(big.Int).MulRange(1, int64(366-numPeople))
		power366 := new(big.Int).Exp(big.NewInt(366), big.NewInt(int64(numPeople)), nil)
		numerator := new(big.Float).SetInt(factorial366)
		denomInt := new(big.Int).Mul(power366, factorialN)
		denom := new(big.Float).SetInt(denomInt)
		prob := new(big.Float).Sub(big.NewFloat(1.0), new(big.Float).Quo(numerator, denom))
		//fmt.Println("Actual probability =", prob)
		fmt.Println("Actual probability =", prob.String())
	}
}
