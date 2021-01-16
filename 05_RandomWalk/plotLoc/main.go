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

func main() {
	steps := []location.Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	var usualDrunk drunk.Drunk
	usualDrunk.SetName("usual")
	usualDrunk.SetStepChoices(steps)

	steps = []location.Location{{0.0, 1.1}, {0.0, -0.9}, {1.1, 0.0}, {-0.9, 0.0}}
	var masochistDrunk drunk.Drunk
	masochistDrunk.SetName("masochist")
	masochistDrunk.SetStepChoices(steps)

	drunks := [...]drunk.Drunk{usualDrunk, masochistDrunk}

	rand.Seed(time.Now().UTC().UnixNano())
	plotLocs(drunks[:], 10000, 1000)
}

func getFinalLocs(numSteps int, numTrials int, dClass drunk.Drunk) []location.Location {
	var locs []location.Location
	for t := 0; t < numTrials; t++ {
		var origin location.Location
		var f field.Field
		f.AddDrunk(dClass, origin)
		for s := 0; s < numSteps; s++ {
			f.MoveDrunk(dClass)
		}
		loc, err := f.GetLoc(dClass)
		if err != nil {
			log.Fatalln("getLoc", err)
			continue
		}
		locs = append(locs, loc)
	}
	return locs
}

func plotLocs(drunkKinds []drunk.Drunk, numSteps int, numTrials int) {
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
			pts[i].X = l.X
			pts[i].Y = l.Y
			sumX += l.X
			sumY += l.Y
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
