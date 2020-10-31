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

func traceWalk(fieldKinds []field.OddField, numSteps int, xRange int, yRange int) {
	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		return
	}

	p.Title.Text = fmt.Sprintf("Spots Visited on Walk (%d) steps", numSteps)
	p.X.Label.Text = "Steps East/West of Origin"
	p.Y.Label.Text = "Steps North/South of Origin"
	p.X.Min = float64(-xRange)
	p.X.Max = float64(xRange)
	p.Y.Min = float64(-yRange)
	p.Y.Max = float64(yRange)
	p.Add(plotter.NewGrid())

	for f, fClass := range fieldKinds {
		var origin location.Location
		steps := []location.Location{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		var usualDrunk drunk.Drunk
		usualDrunk.SetName("usual")
		usualDrunk.SetStepChoices(steps)
		fClass.AddDrunk(usualDrunk, origin)

		var locs []location.Location
		for s := 0; s < numSteps; s++ {
			fClass.MoveDrunk(usualDrunk)
			loc, err := fClass.GetLoc(usualDrunk)
			if err != nil {
				log.Fatalln("getLoc", err)
				continue
			}
			locs = append(locs, loc)
		}

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
		s.GlyphStyle.Color = plotutil.Color(f)
		s.GlyphStyle.Radius = vg.Points(3)

		legend := fmt.Sprintf("%s", fClass.Name())

		p.Add(s)
		p.Legend.Add(legend, s)
	}

	if err := p.Save(8*vg.Inch, 8*vg.Inch, "trace_walk.png"); err != nil {
		log.Fatalln("plot.Save()", err)
		return
	}
}

func main() {
	var fields []field.OddField

	var of field.OddField
	of.SetName("Normal")
	fields = append(fields, of)

	of.SetWormHoles(1000, 100, 100)
	of.SetName("Odd Field")
	fields = append(fields, of)

	rand.Seed(time.Now().UTC().UnixNano())
	traceWalk(fields, 500, 100, 100)
}
