package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gonum/matrix/mat64"
	"github.com/kortschak/nmf"
)

func posNorm(_, _ int, _ float64) float64 { return math.Abs(rand.NormFloat64()) }

func NMF(V *mat64.Dense, k int) (*mat64.Dense, *mat64.Dense, float64, bool) {
	rand.Seed(1)
	categories := k
	rows, cols := V.Dims()
	Wo := mat64.NewDense(rows, categories, nil)
	Wo.Apply(posNorm, Wo)

	Ho := mat64.NewDense(categories, cols, nil)
	Ho.Apply(posNorm, Ho)

	conf := nmf.Config{
		Tolerance:   1e-5,
		MaxIter:     100,
		MaxOuterSub: 1000,
		MaxInnerSub: 20,
		Limit:       time.Second * 100,
	}

	W, H, ok := nmf.Factors(V, Wo, Ho, conf)
	log.Println("W Dims")
	log.Println(W.Dims())
	log.Println("H Dims")
	log.Println(H.Dims())
	var P, D mat64.Dense
	P.Mul(W, H)
	D.Sub(V, &P)
	return W, H, mat64.Norm(&D, 2), ok
}
