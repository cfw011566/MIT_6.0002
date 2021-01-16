package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"../citytemp"
	"gonum.org/v1/gonum/stat"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	population := citytemp.GetHighs()
	popMean := stat.Mean(population, nil)

	sampleSize := 200
	numTrials := 1000

	numBad := 0

	for t := 0; t < numTrials; t++ {
		sample := citytemp.Sampling(sampleSize)
		sampleMean, std := stat.MeanStdDev(sample, nil)
		stdErr := stat.StdErr(std, float64(sampleSize))
		if math.Abs(popMean-sampleMean) > 1.96*stdErr {
			numBad++
		}
	}
	fmt.Printf("Fraction outside 95%% confidence interval = %f\n", float64(numBad)/float64(numTrials))
}
