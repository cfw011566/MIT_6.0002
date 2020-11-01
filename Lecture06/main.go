package main

import (
	"fmt"

	"./roulette"
)

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

func main() {
	//rand.Seed(time.Now().UTC().UnixNano())

	//test_fair()

	test_all()
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
	numSpins := [...]int{1000, 10000, 100000, 1000000}
	var games []roulette.Roulette
	var game roulette.Roulette
	game.Init(roulette.Fair)
	games = append(games, game)
	game.Init(roulette.European)
	games = append(games, game)
	game.Init(roulette.American)
	games = append(games, game)
	for _, spin := range numSpins {
		fmt.Println("\nSimulate", numTrials, "trials of", spin, "spins each")
		for _, game := range games {
			pocketReturns := findPocketReturn(game, numTrials, spin, false)
			sum := 0.0
			for _, r := range pocketReturns {
				sum += r
			}
			expReturn := 100.8 * sum / float64(len(pocketReturns))
			fmt.Printf("Exp. return for %s = %.4f%%\n", game, expReturn)
		}
	}
}
