package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/stat"
)

func throwNeedles(numNeedles int) float64 {
	inCircle := 0
	for i := 0; i < numNeedles; i++ {
		x := rand.Float64()
		y := rand.Float64()
		if math.Sqrt(x*x+y*y) <= 1.0 {
			inCircle++
		}
	}
	return 4.0 * (float64(inCircle) / float64(numNeedles))
}

func getEst(numNeedles int, numTrials int) (float64, float64) {
	var estimates []float64
	for t := 0; t < numTrials; t++ {
		piGuess := throwNeedles(numNeedles)
		estimates = append(estimates, piGuess)
	}
	sDev := stat.StdDev(estimates, nil)
	sum := 0.0
	for _, e := range estimates {
		sum += e
	}
	curEst := sum / float64(len(estimates))

	fmt.Printf("Est. = %f, Std. dev. = %.6f, Needles = %d\n", curEst, sDev, numNeedles)
	return curEst, sDev
}

func estPi(precision float64, numTrials int) float64 {
	numNeedles := 1000
	sDev := precision
	var curEst float64
	for sDev >= precision/2 {
		curEst, sDev = getEst(numNeedles, numTrials)
		numNeedles *= 2
	}
	return curEst
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	//fmt.Println("pi =", throwNeedles(1000000))
	estPi(0.005, 100)
}
