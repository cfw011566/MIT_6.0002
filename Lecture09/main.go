package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	plot_data()

	fit_data()

	mystery_data()

	/*
		x := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		y := []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}

		degree := 2

		quadraticModel(x, y, degree)
	*/

	xVals, yVals := getData("mysteryData.txt")
	degrees := [...]int{2, 4, 8, 12, 16}
	models := genFits(xVals, yVals, degrees[:])
	testFits(models, degrees[:], xVals, yVals, "Mystery Data")
}

func getData(fileName string) (xs []float64, ys []float64) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lineNumber := 0
	input := bufio.NewScanner(f)
	for input.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := input.Text()
		for i, column := range strings.Split(line, " ") {
			f, err := strconv.ParseFloat(column, 64)
			if err != nil {
				log.Fatal(err)
				continue
			}
			if i == 0 {
				ys = append(ys, f)
			}
			if i == 1 {
				xs = append(xs, f)
			}
		}
	}

	return xs, ys
}

func plot_data() {
	masses, distances := getData("springData.txt")
	//fmt.Println(masses)
	//fmt.Println(distances)

	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		panic(err)
	}

	p.Title.Text = "Measured Displacement of Spring"
	p.X.Label.Text = "|Force| (Newtons)"
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Label.Text = "Distance (meters)"
	p.Y.Min = 0.05
	p.Y.Max = 0.5

	//p.Add(plotter.NewGrid())

	points := make(plotter.XYs, len(masses))
	for i, mass := range masses {
		points[i].X = mass * 9.81
		points[i].Y = distances[i]
	}
	s, err := plotter.NewScatter(points)
	if err != nil {
		log.Panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	s.GlyphStyle.Radius = vg.Points(5)
	s.GlyphStyle.Shape = draw.CircleGlyph{}

	p.Add(s)

	if err = p.Save(8*vg.Inch, 8*vg.Inch, "plot.png"); err != nil {
		log.Fatalln("plot.Save()", err)
	}
}

func fit_data() {
	masses, distances := getData("springData2.txt")

	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		panic(err)
	}

	p.Title.Text = "Measured Displacement of Spring"
	p.X.Label.Text = "|Force| (Newtons)"
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Label.Text = "Distance (meters)"
	p.Y.Min = 0.05
	p.Y.Max = 0.5
	p.Legend.Top = true
	p.Legend.Left = true

	xVals := make([]float64, len(masses))
	for i, mass := range masses {
		xVals[i] = mass * 9.81
	}
	yVals := distances

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
	p.Legend.Add("Measured points", s)

	b, a := stat.LinearRegression(xVals, yVals, nil, false)
	fmt.Println("a =", a, "b =", b)
	line := plotter.NewFunction(func(x float64) float64 { return a*x + b })
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = 3
	legend := fmt.Sprintf("Linear fit, k = %.5f", 1/a)
	p.Legend.Add(legend, line)

	p.Add(s, line)

	if err = p.Save(8*vg.Inch, 8*vg.Inch, "fit.png"); err != nil {
		log.Fatalln("plot.Save()", err)
	}
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

func mystery_data() {
	xVals, yVals := getData("mysteryData.txt")

	p, err := plot.New()
	if err != nil {
		log.Fatalln("plot.New()", err)
		panic(err)
	}

	p.Title.Text = "Mystery Data"
	//p.Legend.Top = true
	//p.Legend.Left = true

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
	p.Legend.Add("Data Points", s)

	b, a := stat.LinearRegression(xVals, yVals, nil, false)
	line := plotter.NewFunction(func(x float64) float64 { return a*x + b })
	line.Color = color.RGBA{G: 255, A: 255}
	line.Width = 3
	p.Legend.Add("Liner Model", line)

	model := polyRegression(xVals, yVals, 2)
	quadratic := plotter.NewFunction(func(x float64) float64 { return model[0] + model[1]*x + model[2]*x*x })
	quadratic.Color = color.RGBA{R: 255, A: 255}
	quadratic.Width = 3
	quadratic.Dashes = plotutil.Dashes(1)
	p.Legend.Add("Quadratic Model", quadratic)

	p.Add(s, line, quadratic)

	if err = p.Save(8*vg.Inch, 8*vg.Inch, "mystery.png"); err != nil {
		log.Fatalln("plot.Save()", err)
	}
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

	if err = p.Save(8*vg.Inch, 8*vg.Inch, "mystery2.png"); err != nil {
		log.Fatalln("plot.Save()", err)
	}
}
