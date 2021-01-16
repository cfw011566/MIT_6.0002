package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	numCasesPerYear := 36000
	numYears := 3
	stateSize := 10000
	communitySize := 10
	numCommunities := stateSize / communitySize

	numTrials := 100
	anyRegion := 0
	for t := 0; t < numTrials; t++ {
		locs := make([]int, numCommunities)
		for i := 0; i < numYears*numCasesPerYear; i++ {
			locs[rand.Intn(numCommunities)]++
		}
		max := 0
		for _, n := range locs {
			if n > max {
				max = n
			}
		}
		if max >= 143 {
			anyRegion++
		}
	}
	fmt.Println(anyRegion)
	prob := float64(anyRegion) / float64(numTrials)
	fmt.Printf("Est. probability of some region having at least 143 cases = %.4f\n", prob)
}
