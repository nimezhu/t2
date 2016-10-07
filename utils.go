package main

import (
	"math"
	"sort"
)

func Gini(a []float64) float64 {
	n := len(a)
	b := make([]float64, n)
	copy(b, a)
	sort.Float64s(b)
	s1 := 0.0
	s2 := 0.0
	for i, v := range b {
		s1 += (float64(i) + 1.0) * v
		s2 += v
	}
	if s2 == 0 {
		return 0.0
	}
	n0 := float64(n)
	G := 2.0*s1/(n0*s2) - (n0+1.0)/n0
	return G

}

func EntropyBits(p []float64) float64 {
	e := 0.0
	for _, x := range p {
		if x != 0.0 {
			e += math.Log2(x) * x
		}
	}
	return -e
}
