package main

import (
	"fmt"
	"math"

	"./roulette"
)

func main() {
	//rand.Seed(time.Now().UTC().UnixNano())

	test_fair()

	test_all()

	test_empirical()
}

func playRoulette(game roulette.Roulette, numSpins int, pocket int, bet int, toPrint bool) float64 {
	totPocket := 0
	for i := 0; i < numSpins; i++ {
		game.Spin()
		totPocket += game.BetPocket(pocket, bet)
	}
	expectReturn := float64(totPocket) / float64(numSpins)
	if toPrint {
		fmt.Println(numSpins, "spin of", game)
		fmt.Printf("Expect return betting %d = %.4f%%\n", pocket, 100.0*expectReturn)
	}
	return expectReturn
}

func test_fair() {
	var game roulette.Roulette
	game.Init(roulette.Fair)

	numSpins := [...]int{100, 1000000}
	for _, spin := range numSpins {
		for i := 0; i < 3; i++ {
			playRoulette(game, spin, 2, 1, true)
		}
	}
}

func findPocketReturn(game roulette.Roulette, numTrials int, trialSize int, toPrint bool) []float64 {
	var pocketReturns []float64
	for i := 0; i < numTrials; i++ {
		trivals := playRoulette(game, trialSize, 2, 1, toPrint)
		pocketReturns = append(pocketReturns, trivals)
	}
	return pocketReturns
}

func test_all() {
	numTrials := 20
	var games []roulette.Roulette
	var game roulette.Roulette
	game.Init(roulette.Fair)
	games = append(games, game)
	game.Init(roulette.European)
	games = append(games, game)
	game.Init(roulette.American)
	games = append(games, game)
	for _, numSpins := range [...]int{1000, 10000, 100000, 1000000} {
		fmt.Println("\nSimulate", numTrials, "trials of", numSpins, "spins each")
		for _, game := range games {
			pocketReturns := findPocketReturn(game, numTrials, numSpins, false)
			sum := 0.0
			for _, r := range pocketReturns {
				sum += r
			}
			expReturn := 100.8 * sum / float64(len(pocketReturns))
			fmt.Printf("Exp. return for %s = %.4f%%\n", game, expReturn)
		}
	}
}

func getMeanAndStd(X []float64) (float64, float64) {
	sum := 0.0
	for _, x := range X {
		sum += x
	}
	mean := sum / float64(len(X))
	tot := 0.0
	for _, x := range X {
		tot += (x - mean) * (x - mean)
	}
	std := math.Pow(tot/float64(len(X)), 0.5)
	return mean, std
}

func test_empirical() {
	numTrials := 20
	var games []roulette.Roulette
	var game roulette.Roulette
	game.Init(roulette.Fair)
	games = append(games, game)
	game.Init(roulette.European)
	games = append(games, game)
	game.Init(roulette.American)
	games = append(games, game)
	for _, numSpins := range [...]int{1000, 100000, 1000000} {
		fmt.Println("\nSimulate betting a pocket for", numTrials, "trials of", numSpins, "spin each")
		for _, game := range games {
			pocketReturns := findPocketReturn(game, numTrials, numSpins, false)
			mean, std := getMeanAndStd(pocketReturns)
			fmt.Printf("Exp. return for %s = %.3f%%, +/- %.3f%% with 95%% confidence\n", game, 100.0*mean, 100.0*1.96*std)
		}
	}
}
