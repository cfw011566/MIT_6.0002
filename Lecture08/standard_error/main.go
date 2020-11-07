package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"../citytemp"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	sampleSizes := [...]int{25, 50, 100, 200, 300, 400, 500, 600}
	numTrials := 50

	population := citytemp.GetHighs()
	popSD := stat.StdDev(population, nil)

	var sems []float64
	var sampleSDs []float64

	for _, size := range sampleSizes {
		sem := stat.StdErr(popSD, float64(size))
		sems = append(sems, sem)
		var means []float64
		for t := 0; t < numTrials; t++ {
			sample := citytemp.Sampling(size)
			mean := stat.Mean(sample, nil)
			means = append(means, mean)
		}
		std := stat.StdDev(means, nil)
		sampleSDs = append(sampleSDs, std)
	}
	fmt.Println(sems)
	fmt.Println(sampleSDs)

	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		return
	}

	p.Title.Text = fmt.Sprintf("SD for %d Means and SEM", numTrials)
	p.X.Label.Text = "Sample Size"
	p.Y.Label.Text = "Std and SEM"
	p.Y.Min = 0
	p.Legend.Top = true
	p.Legend.XOffs = -10

	pts := make(plotter.XYs, len(sampleSizes))
	for i, s := range sampleSizes {
		pts[i].X = float64(s)
		pts[i].Y = sampleSDs[i]
	}
	lpLine, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatalln("plot.NewLinePoints()", err)
		return
	}
	lpLine.Color = color.RGBA{B: 255, A: 255}
	lpLine.Dashes = plotutil.Dashes(0)
	lpLine.Width = 3
	p.Add(lpLine)
	legend := fmt.Sprintf("Std of %d means", numTrials)
	p.Legend.Add(legend, lpLine)

	pts = make(plotter.XYs, len(sampleSizes))
	for i, s := range sampleSizes {
		pts[i].X = float64(s)
		pts[i].Y = sems[i]
	}
	lpLine, err = plotter.NewLine(pts)
	if err != nil {
		log.Fatalln("plot.NewLinePoints()", err)
		return
	}
	lpLine.Color = color.RGBA{R: 255, A: 255}
	lpLine.Dashes = plotutil.Dashes(1)
	lpLine.Width = 3
	p.Add(lpLine)
	p.Legend.Add("SEM", lpLine)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "stdErr.png"); err != nil {
		panic(err)
	}
}
