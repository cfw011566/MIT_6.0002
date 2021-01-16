package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"./cluster"
	"gonum.org/v1/gonum/stat"
)

func scaleAttrs(vals []float64) []float64 {
	mean, sd := stat.MeanStdDev(vals, nil)
	for i, v := range vals {
		vals[i] = (v - mean) / sd
	}
	return vals
}

func getData(toScale bool) (patients []cluster.Patient) {
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
		//p := cluster.Patient{Name: n, Features: f[:], Label: classList[i]}
		var p cluster.Patient
		p.Init(n, f[:], classList[i])
		patients = append(patients, p)
	}
	return patients
}

func sampling(patients []cluster.Patient, size int) []cluster.Patient {
	rand.Shuffle(len(patients), func(i, j int) {
		patients[i], patients[j] = patients[j], patients[i]
	})
	samples := patients[:size]
	return samples
}

func kmeans(examples []cluster.Patient, k int, verbose bool) ([]cluster.Cluster, error) {
	initialCentroids := sampling(examples, k)
	//fmt.Println("initialCentroids =", initialCentroids)
	var clusters []cluster.Cluster
	for _, e := range initialCentroids {
		es := [...]cluster.Patient{e}
		c := cluster.Cluster{}
		c.Init(es[:])
		clusters = append(clusters, c)
	}
	//fmt.Println("clusters =", clusters)

	converged := false
	numIterations := 0
	for !converged {
		numIterations += 1
		newClusters := make([][]cluster.Patient, k)
		for i := 0; i < k; i++ {
			newClusters[i] = make([]cluster.Patient, 0)
		}

		for _, e := range examples {
			smallestDistance := e.Distance(clusters[0].Centroid())
			index := 0
			for i := 1; i < k; i++ {
				distance := e.Distance(clusters[i].Centroid())
				if distance < smallestDistance {
					smallestDistance = distance
					index = i
				}
			}
			//fmt.Println(e, "index =", index, "distance =", smallestDistance)
			newClusters[index] = append(newClusters[index], e)
		}
		//fmt.Println("new clusters 0 =", newClusters[0])
		//fmt.Println("new clusters 1 =", newClusters[1])

		for _, c := range newClusters {
			if len(c) == 0 {
				//panic("Empty Cluster")
				err := errors.New("Empty Cluster")
				return clusters, err
			}
		}

		converged = true
		for i := 0; i < k; i++ {
			if clusters[i].Update(newClusters[i]) > 0.0 {
				converged = false
			}
		}

		if verbose {
			fmt.Println("Iteration #", numIterations)
			for _, c := range clusters {
				fmt.Println(c)
			}
			fmt.Println("")
		}
	}
	return clusters, nil
}

func trykmeans(examples []cluster.Patient, numClusters int, numTrials int, verbose bool) []cluster.Cluster {
	best, err := kmeans(examples, numClusters, verbose)
	if err != nil {
		panic(err)
		return best
	}
	minDissimilarity := cluster.Dissimilarity(best)
	trial := 1
	for trial < numTrials {
		clusters, err := kmeans(examples, numClusters, verbose)
		if err != nil {
			continue
		}
		currDissimilarity := cluster.Dissimilarity(clusters)
		if currDissimilarity < minDissimilarity {
			copy(best, clusters)
			minDissimilarity = currDissimilarity
		}
		trial++
	}
	return best
}

func printClustering(clustering []cluster.Cluster) []float64 {
	var posFracs []float64
	for _, c := range clustering {
		numPts := 0
		numPos := 0
		for _, p := range c.Members() {
			numPts++
			if p.Label() > 0.5 {
				numPos++
			}
		}
		fracPos := float64(numPos) / float64(numPts)
		posFracs = append(posFracs, fracPos)
		fmt.Printf("Cluster of size %d with fraction of positives = %.4f\n", numPts, fracPos)
	}
	return posFracs
}

func testClustering(patients []cluster.Patient, numClusters int, numTrials int) []float64 {
	bestClustering := trykmeans(patients, numClusters, numTrials, false)
	posFracs := printClustering(bestClustering)
	return posFracs
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	//patients := getData(false)
	//fmt.Println(patients)
	patients := getData(true)
	//fmt.Println(patients)

	numClusters := [...]int{2, 4, 6}
	for _, k := range numClusters {
		fmt.Printf("\nTest k-means (k = %d)\n", k)
		testClustering(patients, k, 2)
	}
}
