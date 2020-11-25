package cluster

import (
	"fmt"
	"math"
	"sort"
)

func minkowskiDist(p1, p2 Patient, p int) float64 {
	v1 := p1.Features()
	v2 := p2.Features()
	dist := 0.0
	for i := 0; i < len(v1); i++ {
		dist += math.Pow((math.Abs(v1[i] - v2[i])), float64(p))
	}
	return math.Pow(dist, 1/float64(p))
}

type Patient struct {
	name     string
	features []float64
	label    float64
}

func (p *Patient) Init(name string, features []float64, label float64) {
	p.name = name
	p.features = make([]float64, len(features))
	copy(p.features, features)
	p.label = label
}

func (p Patient) Dimensionality() int { return len(p.features) }

func (p Patient) Features() []float64 { return p.features }

func (p Patient) Label() float64 { return p.label }

func (p Patient) Name() string { return p.name }

func (p Patient) Distance(q Patient) float64 { return minkowskiDist(p, q, 2) }

func (p Patient) String() string {
	fString := ""
	for i, f := range p.features {
		if i == 0 {
			fString += fmt.Sprintf("%.4f", f)
		} else {
			fString += fmt.Sprintf(", %.4f", f)
		}
	}
	return fmt.Sprintf("%s:%s:%.4f", p.name, fString, p.label)
}

type Cluster struct {
	examples []Patient
	centroid Patient
}

func (c *Cluster) Init(examples []Patient) {
	c.examples = make([]Patient, len(examples))
	copy(c.examples, examples)
	//fmt.Println("c len =", len(c.examples), "in len =", len(examples))
	c.centroid = c.ComputeCentroid()
}

func (c *Cluster) Update(examples []Patient) float64 {
	oldCentroid := c.centroid
	//fmt.Println("c.centroid =", c.centroid)
	//fmt.Println("oldCentroid =", oldCentroid)
	c.examples = make([]Patient, len(examples))
	copy(c.examples, examples)
	//fmt.Println("examples =", examples)
	//fmt.Println("c.examples =", c.examples)
	c.centroid = c.ComputeCentroid()
	return oldCentroid.Distance(c.centroid)
}

func (c *Cluster) ComputeCentroid() Patient {
	//fmt.Println(c)
	vals := make([]float64, c.examples[0].Dimensionality())
	for _, e := range c.examples {
		for i, f := range e.features {
			vals[i] += f
		}
	}
	means := make([]float64, len(vals))
	for i := 0; i < len(vals); i++ {
		means[i] = vals[i] / float64(len(c.examples))
	}
	centroid := Patient{"centroid", means, 0.0}
	return centroid
}

func (c Cluster) Centroid() Patient { return c.centroid }

func (c Cluster) Variability() float64 {
	totDist := 0.0
	for _, e := range c.examples {
		totDist += math.Pow(e.Distance(c.centroid), 2.0)
	}
	return totDist
}

func (c Cluster) Members() []Patient {
	return c.examples
}

func (c Cluster) String() string {
	var names []string
	for _, e := range c.examples {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	result := "Cluster with centroid "
	for i, f := range c.centroid.Features() {
		if i == 0 {
			result += fmt.Sprintf("%.4f", f)
		} else {
			result += fmt.Sprintf(", %.4f", f)
		}
	}
	result += " contains:\n"
	for i, n := range names {
		if i == 0 {
			result += n
		} else {
			result += ", " + n
		}
	}
	return result
}

func Dissimilarity(clusters []Cluster) float64 {
	totDist := 0.0
	for _, c := range clusters {
		totDist += c.Variability()
	}
	return totDist
}
