package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"../drunk"
	"../field"
	"../location"
)

// functions
func walk(f field.Field, d drunk.Drunk, numSteps int) float64 {
	start, err := f.GetLoc(d)
	if err != nil {
		log.Fatalln("error", err)
		return 0.0
	}
	for s := 0; s < numSteps; s++ {
		f.MoveDrunk(d)
	}
	loc, err := f.GetLoc(d)
	if err != nil {
		log.Fatalln("error", err)
		return 0.0
	}
	return start.DistFrom(loc)
}

func simWalks(numSteps int, numTrials int, dClass drunk.Drunk) []float64 {
	var origin location.Location
	var distances []float64
	for t := 0; t < numTrials; t++ {
		var f field.Field
		f.AddDrunk(dClass, origin)
		distances = append(distances, walk(f, dClass, numSteps))
	}
	return distances
}

func drunkTest(walkLengths []int, numTrials int, dClass drunk.Drunk) {
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

func simDrunk(numTrials int, dClass drunk.Drunk, walkLengths []int) []float64 {
	var meanDistances []float64
	for _, numSteps := range walkLengths {
		fmt.Println("Start simulation of", numSteps, "steps")
		trials := simWalks(numSteps, numTrials, dClass)
		sum := 0.0
		for _, d := range trials {
			sum += d
		}
		mean := sum / float64(len(trials))
		meanDistances = append(meanDistances, mean)
	}

	return meanDistances
}

func simAll(drunkKinds []drunk.Drunk, walkLengths []int, numTrials int) {
	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		return
	}

	p.Title.Text = fmt.Sprintf("Mean Distance from Origin (%d) trials", numTrials)
	p.X.Label.Text = "Number of Steps"
	p.Y.Label.Text = "Distance from Origin"
	p.Add(plotter.NewGrid())

	for d, dClass := range drunkKinds {
		fmt.Println("Start simulation of", dClass)
		means := simDrunk(numTrials, dClass, walkLengths)
		fmt.Println("means =", means)
		pts := make(plotter.XYs, len(walkLengths))
		for i, w := range walkLengths {
			pts[i].X = float64(w)
			pts[i].Y = means[i]
		}
		lpLine, lpPoints, err := plotter.NewLinePoints(pts)
		if err != nil {
			log.Fatalln("plot.NewLinePoints()", err)
			continue
		}
		lpLine.Color = plotutil.Color(d)
		lpLine.Dashes = plotutil.Dashes(d)
		lpPoints.Shape = plotutil.Shape(d)
		lpPoints.Color = plotutil.Color(d)

		p.Add(lpPoints, lpLine)
		p.Legend.Add(dClass.Name(), lpLine, lpPoints)

		//if err = plotutil.AddLinePoints(p, dClass.name, pts); err != nil {
		//	log.Fatalln("AddLinePoints", err)
		//}
	}

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
		log.Fatalln("plot.Save()", err)
		return
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	test_sanity()

	test_walk()

	test_plot_all()
}

func test_sanity() {
	p := location.Location{1.2, 2.3}
	fmt.Println("p=", p)

	steps := []location.Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	var usualDrunk drunk.Drunk
	usualDrunk.SetName("usual")
	usualDrunk.SetStepChoices(steps)
	fmt.Println(usualDrunk)

	steps = []location.Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	var masochistDrunk drunk.Drunk
	masochistDrunk.SetName("masochist")
	masochistDrunk.SetStepChoices(steps)
	fmt.Println(masochistDrunk)

	var f field.Field
	var origin location.Location
	fmt.Println("Field", f)
	f.AddDrunk(usualDrunk, origin)
	fmt.Println("add usual", f)
	f.AddDrunk(masochistDrunk, origin)
	fmt.Println("add masochist", f)
	dist := walk(f, usualDrunk, 10000)
	fmt.Println("distance=", dist)
	dist = walk(f, masochistDrunk, 10000)
	fmt.Println("distance=", dist)
}

func test_walk() {
	steps := []location.Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	var usualDrunk drunk.Drunk
	usualDrunk.SetName("usual")
	usualDrunk.SetStepChoices(steps)

	steps = []location.Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	var masochistDrunk drunk.Drunk
	masochistDrunk.SetName("masochist")
	masochistDrunk.SetStepChoices(steps)

	testSteps := [...]int{1000, 10000}
	drunkTest(testSteps[:], 100, usualDrunk)
	drunkTest(testSteps[:], 100, masochistDrunk)
}

func test_plot_all() {
	steps := []location.Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	var usualDrunk drunk.Drunk
	usualDrunk.SetName("usual")
	usualDrunk.SetStepChoices(steps)

	steps = []location.Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	var masochistDrunk drunk.Drunk
	masochistDrunk.SetName("masochist")
	masochistDrunk.SetStepChoices(steps)

	drunks := [...]drunk.Drunk{usualDrunk, masochistDrunk}
	numSteps := [...]int{10, 100, 1000, 10000, 100000}
	simAll(drunks[:], numSteps[:], 100)
}
