package main

import (
	"fmt"
	"math/rand"
	"time"

	"../citytemp"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func makeHist(data []float64, fileName, title, xLabel, yLabel string, bins int) {
	v := make(plotter.Values, len(data))
	for i := range v {
		v[i] = data[i]
	}
	// Make a plot and set its title.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Create a histogram of our values drawn
	// from the standard normal.
	h, err := plotter.NewHist(v, bins)
	if err != nil {
		panic(err)
	}
	h.FillColor = plotutil.Color(rand.Intn(5))
	//h.FillColor = color.RGBA{0, 0, 255, 255}
	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, fileName); err != nil {
		panic(err)
	}

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	population := citytemp.GetHighs()

	mean, std := stat.MeanStdDev(population, nil)
	fmt.Printf("Population mean = %.2f\n", mean)
	fmt.Printf("Standard deviation of population = %.2f\n", std)

	title := fmt.Sprintf("Daily High 1961-2015, Population\n(mean = %.2f)", mean)
	xLabel := "Degrees C"
	yLabel := "Number Days"
	makeHist(population, "population.png", title, xLabel, yLabel, 20)

	samples := citytemp.Sampling(100)
	mean, std = stat.MeanStdDev(samples, nil)
	fmt.Printf("Sample mean = %.2f\n", mean)
	fmt.Printf("Standard deviation of sample = %.2f\n", std)

	title = fmt.Sprintf("Daily High 1961-2015, Sample\n(mean = %.2f)", mean)
	xLabel = "Degrees C"
	yLabel = "Number Days"
	makeHist(samples, "sample100.png", title, xLabel, yLabel, 20)

	sampleSize := 100
	numSamples := 1000
	var sampleMeans []float64
	for i := 0; i < numSamples; i++ {
		sample := citytemp.Sampling(sampleSize)
		mean := stat.Mean(sample, nil)
		sampleMeans = append(sampleMeans, mean)
	}
	mean, std = stat.MeanStdDev(sampleMeans, nil)
	fmt.Printf("Mean of sample means = %.2f\n", mean)
	fmt.Printf("Standard deviation of sample means = %.2f\n", std)
	makeHist(sampleMeans, "sampleAll.png", "Means of Samples", "Mean", "Frequency", 20)
}
