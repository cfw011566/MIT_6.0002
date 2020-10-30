package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

func (d Drunk) Name() string { return d.name }
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

func simDrunk(numTrials int, dClass Drunk, walkLengths []int) []float64 {
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

func simAll(drunkKinds []Drunk, walkLengths []int, numTrials int) {
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
		p.Legend.Add(dClass.name, lpLine, lpPoints)

		//if err = plotutil.AddLinePoints(p, dClass.name, pts); err != nil {
		//	log.Fatalln("AddLinePoints", err)
		//}
	}

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
		log.Fatalln("plot.Save()", err)
		return
	}
}

func getFinalLocs(numSteps int, numTrials int, dClass Drunk) []Location {
	var locs []Location
	for t := 0; t < numTrials; t++ {
		f := Field{map[string]Location{}}
		f.addDrunk(dClass, Location{0.0, 0.0})
		for s := 0; s < numSteps; s++ {
			f.moveDrunk(dClass)
		}
		loc, err := f.getLoc(dClass)
		if err != nil {
			log.Fatalln("getLoc", err)
			continue
		}
		locs = append(locs, loc)
	}
	return locs
}

func plotLocs(drunkKinds []Drunk, numSteps int, numTrials int) {
	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		return
	}

	p.Title.Text = fmt.Sprintf("Location at End of Walks (%d steps)", numSteps)
	p.X.Label.Text = "Steps East/West of Origin"
	p.Y.Label.Text = "Steps North/South of Origin"
	p.X.Min = -1000
	p.X.Max = 1000
	p.Y.Min = -1000
	p.Y.Max = 1000
	p.Add(plotter.NewGrid())

	for d, dClass := range drunkKinds {
		locs := getFinalLocs(numSteps, numTrials, dClass)
		pts := make(plotter.XYs, len(locs))
		sumX := 0.0
		sumY := 0.0
		for i, l := range locs {
			pts[i].X = l.X()
			pts[i].Y = l.Y()
			sumX += l.X()
			sumY += l.Y()
		}
		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Panic(err)
		}
		s.GlyphStyle.Color = plotutil.Color(d)
		s.GlyphStyle.Radius = vg.Points(3)

		meanX := sumX / float64(len(locs))
		meanY := sumY / float64(len(locs))
		legend := fmt.Sprintf("%s mean abs dist = <%.4f, %.4f>", dClass.Name(), meanX, meanY)

		p.Add(s)
		p.Legend.Add(legend, s)
	}

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "scatter.png"); err != nil {
		log.Fatalln("plot.Save()", err)
		return
	}
}

func main() {
	//rand.Seed(time.Now().UTC().UnixNano())

	//test_sanity()

	//test_walk()

	//test_plot_all()

	test_plot_loc()
}

func test_sanity() {
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

func test_wak() {
	steps := []Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	usualDrunk := Drunk{"usual", steps}

	steps = []Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	masochistDrunk := Drunk{"masochist", steps}

	testSteps := [...]int{1000, 10000}
	drunkTest(testSteps[:], 100, usualDrunk)
	drunkTest(testSteps[:], 100, masochistDrunk)
}

func test_plot_all() {
	steps := []Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	usualDrunk := Drunk{"usual", steps}

	steps = []Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	masochistDrunk := Drunk{"masochist", steps}

	drunks := [...]Drunk{usualDrunk, masochistDrunk}
	numSteps := [...]int{10, 100, 1000, 10000, 100000}
	simAll(drunks[:], numSteps[:], 100)
}

func test_plot_loc() {
	steps := []Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	usualDrunk := Drunk{"usual", steps}

	steps = []Location{{0.0, 1.1}, {0.0, -0.9}, {1.0, 0.0}, {-1.0, 0.0}}
	masochistDrunk := Drunk{"masochist", steps}

	drunks := [...]Drunk{usualDrunk, masochistDrunk}
	plotLocs(drunks[:], 10000, 1000)
}
