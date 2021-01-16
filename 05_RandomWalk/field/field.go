package field

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	"../drunk"
	"../location"
)

type Field struct {
	drunks map[string]location.Location // key: name of drunk
}

func (f *Field) AddDrunk(drunk drunk.Drunk, loc location.Location) {
	if f.drunks == nil {
		f.drunks = make(map[string]location.Location)
	}
	_, ok := f.drunks[drunk.Name()]
	if ok {
		log.Fatalln("Duplicate Drunk")
	} else {
		f.drunks[drunk.Name()] = loc
	}
}

func (f *Field) GetLoc(drunk drunk.Drunk) (location.Location, error) {
	loc, ok := f.drunks[drunk.Name()]
	if ok {
		return loc, nil
	} else {
		var origin location.Location
		return origin, errors.New(fmt.Sprintf("getLoc: Drunk %s not in the field", drunk.Name()))
	}
}

func (f *Field) MoveDrunk(drunk drunk.Drunk) error {
	loc, ok := f.drunks[drunk.Name()]
	if !ok {
		return errors.New("moveDrunk: Drunk not in the field")
	}
	xDist, yDist := drunk.TakeStep()
	f.drunks[drunk.Name()] = loc.Move(xDist, yDist)

	return nil
}

// OddField
type OddField struct {
	Field
	name      string
	wormHoles map[location.Location]location.Location
}

func (f *OddField) Name() string        { return f.name }
func (f *OddField) SetName(name string) { f.name = name }
func (f *OddField) SetWormHoles(numHoles int, xRange int, yRange int) {
	f.wormHoles = map[location.Location]location.Location{}
	for w := 0; w < numHoles; w++ {
		x := float64(rand.Intn(2*xRange) - xRange)
		y := float64(rand.Intn(2*yRange) - yRange)
		loc := location.Location{x, y}
		newX := float64(rand.Intn(2*xRange) - xRange)
		newY := float64(rand.Intn(2*yRange) - yRange)
		newLoc := location.Location{newX, newY}
		f.wormHoles[loc] = newLoc
	}
}

func (f *OddField) MoveDrunk(drunk drunk.Drunk) error {
	loc, ok := f.drunks[drunk.Name()]
	if !ok {
		return errors.New("OddField moveDrunk: Drunk not in the field")
	}
	xDist, yDist := drunk.TakeStep()
	nextLoc := loc.Move(xDist, yDist)
	newLoc, ok := f.wormHoles[nextLoc]
	if !ok {
		f.drunks[drunk.Name()] = nextLoc
	} else {
		f.drunks[drunk.Name()] = newLoc
	}

	return nil
}
