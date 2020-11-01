package roulette

import (
	"math/rand"
)

type RouletteType int

const (
	Fair RouletteType = iota
	European
	American
)

type Roulette struct {
	rouletteType RouletteType
	pocketOdd    int
	pockets      []int
	ball         int
}

func (r Roulette) String() string {
	switch r.rouletteType {
	case Fair:
		return "Fair Roulette"
	case European:
		return "European Roulette"
	case American:
		return "American Roulette"
	default:
		return "Fair Roulette"
	}
}

func (r *Roulette) Init(rouletteType RouletteType) {
	r.pockets = make([]int, 0)
	for i := 1; i < 37; i++ {
		r.pockets = append(r.pockets, i)
	}
	r.pocketOdd = len(r.pockets) - 1
	if rouletteType == European {
		r.pockets = append(r.pockets, 0)
	}
	if rouletteType == American {
		r.pockets = append(r.pockets, 0)
		r.pockets = append(r.pockets, 0)
	}
}

func (r *Roulette) Spin() {
	i := rand.Intn(len(r.pockets))
	r.ball = r.pockets[i]
}

func (r *Roulette) BetPocket(pocket int, amt int) int {
	if pocket == r.ball {
		return amt * r.pocketOdd
	} else {
		return -amt
	}
}
