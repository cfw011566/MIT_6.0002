package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
)

type Patient struct {
	name     string
	features []float64
	label    float64
}

func scaleAttrs(vals []float64) []float64 {
	mean, sd := stat.MeanStdDev(vals, nil)
	for i, v := range vals {
		vals[i] = (v - mean) / sd
	}
	return vals
}

func getData(toScale bool) (patients []Patient) {
	f, err := os.Open("cardiacData.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var hrList, stElevList, ageList, prevACSList, classList []float64
	for _, r := range records {
		hr, err := strconv.ParseFloat(r[0], 64)
		stElev, err := strconv.ParseFloat(r[1], 64)
		age, err := strconv.ParseFloat(r[2], 64)
		prevACS, err := strconv.ParseFloat(r[3], 64)
		class, err := strconv.ParseFloat(r[4], 64)
		if err != nil {
			continue
		}
		hrList = append(hrList, hr)
		stElevList = append(stElevList, stElev)
		ageList = append(ageList, age)
		prevACSList = append(prevACSList, prevACS)
		classList = append(classList, class)
	}
	if toScale {
		hrList = scaleAttrs(hrList)
		stElevList = scaleAttrs(stElevList)
		ageList = scaleAttrs(ageList)
		prevACSList = scaleAttrs(prevACSList)
	}
	for i := 0; i < len(hrList); i++ {
		n := fmt.Sprintf("P%03d", i)
		f := [...]float64{hrList[i], stElevList[i], ageList[i], prevACSList[i]}
		p := Patient{name: n, features: f[:], label: classList[i]}
		patients = append(patients, p)
	}
	return patients
}

func main() {
	patients := getData(false)
	fmt.Println(patients)
	patients = getData(true)
	fmt.Println(patients)
}
