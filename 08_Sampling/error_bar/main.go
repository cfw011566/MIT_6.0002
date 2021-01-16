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
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type errPoints struct {
	plotter.XYs
	plotter.YErrors
}

func showErrorBars(population []float64, sizes []int, numTrials int) {
	var xVals []int
	var sizeMeans, sizeSDs []float64
	for _, sampleSize := range sizes {
		xVals = append(xVals, sampleSize)
		var trialMeans []float64
		for t := 0; t < numTrials; t++ {
			sample := citytemp.Sampling(sampleSize)
			sampleMean := stat.Mean(sample, nil)
			trialMeans = append(trialMeans, sampleMean)
		}
		mean, std := stat.MeanStdDev(trialMeans, nil)
		sizeMeans = append(sizeMeans, mean)
		sizeSDs = append(sizeSDs, std)
	}
	//fmt.Println(sizeMeans)
	fmt.Println(sizeSDs)

	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		return
	}

	p.Title.Text = fmt.Sprintf("Mean Temperature (%d trials)", numTrials)
	p.X.Label.Text = "Sample Size"
	p.Y.Label.Text = "Mean"
	p.X.Min = 0
	p.X.Max = float64(sizes[len(sizes)-1] + 10)

	pts := make(plotter.XYs, len(xVals))
	yerrors := make(plotter.YErrors, len(xVals))
	for i, x := range xVals {
		pts[i].X = float64(x)
		pts[i].Y = sizeMeans[i]
		std := sizeSDs[i]
		yerrors[i].Low = -1.96 * std
		yerrors[i].High = 1.96 * std
	}

	data := errPoints{
		XYs:     pts,
		YErrors: yerrors,
	}
	s, err := plotter.NewScatter(data)
	if err != nil {
		log.Panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Points(5)
	s.GlyphStyle.Shape = draw.CircleGlyph{}

	p.Add(s)

	yerrs, err := plotter.NewYErrorBars(data)
	if err != nil {
		log.Panic(err)
	}
	p.Add(yerrs)

	p.Legend.Add("95% Confidence Interval", s)
	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "errorBar.png"); err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	population := citytemp.GetHighs()

	sampleSizes := [...]int{50, 100, 200, 300, 400, 500, 600}
	showErrorBars(population, sampleSizes[:], 100)
}
