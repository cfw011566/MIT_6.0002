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
	xVals1, yVals1 := getData("Dataset 1.txt")
	degrees := [...]int{2, 4, 8, 12, 16}
	models1 := genFits(xVals1, yVals1, degrees[:])
	testFits(models1, degrees[:], xVals1, yVals1, "Dataset 1")

	xVals2, yVals2 := getData("Dataset 2.txt")
	models2 := genFits(xVals2, yVals2, degrees[:])
	testFits(models2, degrees[:], xVals2, yVals2, "Dataset 2")

	testFits(models1, degrees[:], xVals2, yVals2, "DataSet 2-Model 1")
	testFits(models2, degrees[:], xVals1, yVals1, "DataSet 1-Model 2")
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
		fmt.Println("title=", title)
		log.Fatalln("plot.Save()", err)
	}
}
