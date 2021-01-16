package citytemp

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type CityTemp struct {
	city string
	date time.Time
	temp float64
}

var temps []CityTemp

var population []float64

func init() {
	f, err := os.Open("../temperatures.csv")
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
		f, err := strconv.ParseFloat(r[1], 64)
		if err == nil {
			population = append(population, f)
		}
	}
}

func GetHighs() []float64 {
	return population
}

func Sampling(size int) []float64 {
	rand.Shuffle(len(population), func(i, j int) {
		population[i], population[j] = population[j], population[i]
	})
	samples := population[:size]
	return samples
}
