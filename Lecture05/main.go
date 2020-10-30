package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
)

// Location
type Location struct {
	x float64
	y float64
}

func (l Location) X() float64 { return l.x }
func (l Location) Y() float64 { return l.y }
func (l Location) String() string {
	return fmt.Sprintf("<%f,%f>", l.x, l.y)
}

func (l Location) move(deltaX, deltaY float64) Location {
	return Location{l.x + deltaX, l.y + deltaY}
}
func (l Location) distFrom(other Location) float64 {
	return math.Hypot(l.x-other.x, l.y-other.y)
}

// Drunk
type Drunk struct {
	name        string
	stepChoices []Location
}

func (d Drunk) String() string {
	return fmt.Sprintf("name=%q, steps=%v", d.name, d.stepChoices)
}

func (d Drunk) takeStep() (float64, float64) {
	n := rand.Intn(len(d.stepChoices))
	step := d.stepChoices[n]
	return step.x, step.y
}

func (d Drunk) setStepChoices(steps []Location) {
	d.stepChoices = steps[:]
}

// Field
type Field struct {
	drunks map[string]Location // key: name of drunk
}

func (f Field) addDrunk(drunk Drunk, loc Location) {
	_, ok := f.drunks[drunk.name]
	if ok {
		log.Fatalln("Duplicate Drunk")
	} else {
		f.drunks[drunk.name] = loc
	}
}

func (f Field) getLoc(drunk Drunk) (Location, error) {
	loc, ok := f.drunks[drunk.name]
	if ok {
		return loc, nil
	} else {
		return Location{0.0, 0.0}, errors.New(fmt.Sprintf("getLoc: Drunk %s not in the field", drunk.name))
	}
}

func (f Field) moveDrunk(drunk Drunk) error {
	loc, ok := f.drunks[drunk.name]
	if !ok {
		return errors.New("moveDrunk: Drunk not in the field")
	}
	xDist, yDist := drunk.takeStep()
	f.drunks[drunk.name] = loc.move(xDist, yDist)

	return nil
}

// functions
func walk(f Field, d Drunk, numSteps int) float64 {
	start, err := f.getLoc(d)
	if err != nil {
		log.Fatalln("error", err)
		return 0.0
	}
	for s := 0; s < numSteps; s++ {
		f.moveDrunk(d)
	}
	loc, err := f.getLoc(d)
	if err != nil {
		log.Fatalln("error", err)
		return 0.0
	}
	return start.distFrom(loc)
}

func simWalks(numSteps int, numTrials int, dClass Drunk) []float64 {
	origin := Location{0.0, 0.0}
	var distances []float64
	for t := 0; t < numTrials; t++ {
		var f Field
		f.drunks = map[string]Location{}
		f.addDrunk(dClass, origin)
		distances = append(distances, walk(f, dClass, numSteps))
	}
	return distances
}

func drunkTest(walkLengths []int, numTrials int, dClass Drunk) {
	for _, numSteps := range walkLengths {
		distances := simWalks(numSteps, numTrials, dClass)
		fmt.Println(dClass, "random walk of", numSteps, "steps")
		sum := 0.0
		min := 0.0
		max := 0.0
		for i, d := range distances {
			sum += d
			if i == 0 {
				min = d
				max = d
			} else {
				if d > max {
					max = d
				}
				if d < min {
					min = d
				}
			}
		}
		fmt.Println(" Mean =", sum/float64(len(distances)))
		fmt.Println(" Max =", max, " Min =", min)
	}
}

func main() {
	//rand.Seed(time.Now().UTC().UnixNano())
	test()

	steps := []Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	usualDrunk := Drunk{"usual", steps}

	steps = []Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	masochistDrunk := Drunk{"masochist", steps}

	testSteps := [...]int{1000, 10000}
	drunkTest(testSteps[:], 100, usualDrunk)
	drunkTest(testSteps[:], 100, masochistDrunk)
}

func test() {
	p := Location{1.2, 2.3}
	fmt.Println("p=", p)

	steps := []Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	usualDrunk := Drunk{"usual", steps}
	fmt.Println(usualDrunk)

	steps = []Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	masochistDrunk := Drunk{"masochist", steps}
	fmt.Println(masochistDrunk)

	var f Field
	f.drunks = map[string]Location{}
	//f := Field{map[string]Location{}}
	fmt.Println("Field", f)
	f.addDrunk(usualDrunk, Location{0.0, 0.0})
	fmt.Println("add usual", f)
	f.addDrunk(masochistDrunk, Location{0.0, 0.0})
	fmt.Println("add masochist", f)
	dist := walk(f, usualDrunk, 10000)
	fmt.Println("distance=", dist)
	dist = walk(f, masochistDrunk, 10000)
	fmt.Println("distance=", dist)
}
