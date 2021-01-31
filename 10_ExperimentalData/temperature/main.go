package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type tempDatum struct {
	high float64
	year int
}

func getTempData() (data []tempDatum) {
	f, err := os.Open("temperatures.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, r := range records {
		if i == 0 {
			continue
		}
		h, err := strconv.ParseFloat(r[1], 64)
		if err != nil {
			continue
		}
		y, err := strconv.Atoi(r[2][:4])
		if err != nil {
			continue
		}
		d := tempDatum{high: h, year: y}
		data = append(data, d)
	}
	return data
}

func getYearlyMeans(data []tempDatum) (years []tempDatum) {
	yearData := make(map[int][]float64)
	for _, d := range data {
		y := d.year
		highs := yearData[y]
		yearData[y] = append(highs, d.high)
	}
	for y, highs := range yearData {
		years = append(years, tempDatum{high: stat.Mean(highs, nil), year: y})
	}
	return years
}

func plot_data() {
	data := getTempData()
	years := getYearlyMeans(data)
	sort.Slice(years, func(i, j int) bool {
		return years[i].year < years[j].year
	})

	var xVals, yVals []float64
	for _, y := range years {
		xVals = append(xVals, float64(y.year))
		yVals = append(yVals, y.high)
	}

	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		panic(err)
	}

	p.Title.Text = "Select U.S. Cities"
	p.X.Label.Text = "Year"
	p.Y.Label.Text = "Mean Daily High (C)"
	p.Y.Min = 15.0
	p.Y.Max = 18.0

	points := make(plotter.XYs, len(xVals))
	for i, x := range xVals {
		points[i].X = x
		points[i].Y = yVals[i]
	}
	line, err := plotter.NewLine(points)
	if err != nil {
		log.Fatalln("plot.NewLine()", err)
		return
	}
	line.Color = color.RGBA{B: 255, A: 255}
	line.Dashes = plotutil.Dashes(0)
	line.Width = 3
	p.Add(line)

	if err = p.Save(8*vg.Inch, 8*vg.Inch, "yearly_mean.png"); err != nil {
		log.Fatalln("plot.Save()", err)
	}
}

func splitData(xVals, yVals []float64) (trainX, trainY, testX, testY []float64) {
	var xs []int
	for i := 0; i < len(xVals); i++ {
		xs = append(xs, i)
	}
	rand.Shuffle(len(xs), func(i, j int) {
		xs[i], xs[j] = xs[j], xs[i]
	})
	toTrain := make(map[int]bool)
	for i, x := range xs {
		if i >= len(xs)/2 {
			break
		}
		toTrain[x] = true
	}
	for i := 0; i < len(xVals); i++ {
		_, ok := toTrain[i]
		if ok {
			trainX = append(trainX, xVals[i])
			trainY = append(trainY, yVals[i])
		} else {
			testX = append(testX, xVals[i])
			testY = append(testY, yVals[i])
		}
	}
	return trainX, trainY, testX, testY
}

func train_and_test() {
	data := getTempData()
	years := getYearlyMeans(data)
	sort.Slice(years, func(i, j int) bool {
		return years[i].year < years[j].year
	})

	var xVals, yVals []float64
	for _, y := range years {
		xVals = append(xVals, float64(y.year))
		yVals = append(yVals, y.high)
	}

	numSubsets := 10
	dimensions := [...]int{1, 2, 3, 4}
	rSquares := make(map[int][]float64)
	for _, d := range dimensions {
		rSquares[d] = []float64{}
	}

	for f := 0; f < numSubsets; f++ {
		trainX, trainY, testX, testY := splitData(xVals, yVals)
		for _, d := range dimensions {
			model := polyRegression(trainX, trainY, d)
			var estYVals []float64
			for _, x := range testX {
				y := 0.0
				for e, c := range model {
					y += c * math.Pow(x, float64(e))
				}
				estYVals = append(estYVals, y)
			}
			rSquares[d] = append(rSquares[d], stat.RSquaredFrom(estYVals, testY, nil))
		}
	}
	fmt.Println("Mean R-squares for test data")
	for _, d := range dimensions {
		if d == 1 {
			fmt.Println(rSquares[d])
		}
		mean, std := stat.MeanStdDev(rSquares[d], nil)
		fmt.Printf("For dimensionality %d mean = %.4f Std = %.4f\n", d, mean, std)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	plot_data()

	train_and_test()
}

func Vandermonde(a []float64, degree int) *mat.Dense {
	x := mat.NewDense(len(a), degree+1, nil)
	for i := range a {
		for j, p := 0, 1.; j <= degree; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}

func polyRegression(xs []float64, ys []float64, degree int) []float64 {
	a := Vandermonde(xs, degree)
	b := mat.NewDense(len(ys), 1, ys)
	c := mat.NewDense(degree+1, 1, nil)

	qr := new(mat.QR)
	qr.Factorize(a)

	err := qr.SolveTo(c, false, b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.3f\n", mat.Formatted(c))
	}
	model := mat.Col(nil, 0, c)
	return model
}

func genFits(xVals []float64, yVals []float64, degrees []int) [][]float64 {
	var models [][]float64
	for _, d := range degrees {
		model := polyRegression(xVals, yVals, d)
		models = append(models, model)
	}
	return models
}

func testFits(models [][]float64, degrees []int, xVals []float64, yVals []float64, title string) {
	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		panic(err)
	}

	p.Title.Text = title

	points := make(plotter.XYs, len(xVals))
	for i, x := range xVals {
		points[i].X = x
		points[i].Y = yVals[i]
	}

	s, err := plotter.NewScatter(points)
	if err != nil {
		log.Panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	s.GlyphStyle.Radius = vg.Points(5)
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Legend.Add("Data", s)
	p.Add(s)

	for i, d := range degrees {
		model := models[i]
		line := plotter.NewFunction(func(x float64) float64 {
			sum := 0.0
			for e, c := range model {
				sum += c * math.Pow(x, float64(e))
			}
			return sum
		})
		line.Color = plotutil.Color(i)
		line.Width = 3
		var estYVals []float64
		for _, x := range xVals {
			y := 0.0
			for e, c := range model {
				y += c * math.Pow(x, float64(e))
			}
			estYVals = append(estYVals, y)
		}
		r2 := stat.RSquaredFrom(estYVals, yVals, nil)
		legend := fmt.Sprintf("Fit of degree %d, R2 = %.5f", d, r2)
		p.Legend.Add(legend, line)
		p.Add(line)
	}

	if err = p.Save(8*vg.Inch, 8*vg.Inch, title+".png"); err != nil {
		log.Fatalln("plot.Save()", err)
	}
}
