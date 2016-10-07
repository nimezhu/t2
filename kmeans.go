package main

import (
	//"log"
	"math"
	"sort"

	"github.com/gonum/matrix/mat64"
	"github.com/montanaflynn/stats"
	. "github.com/soniakeys/cluster"
)

func rowToPoints(m mat64.Dense) []Point {
	r, _ := m.Dims()
	points := make([]Point, r)
	for i := range points {
		points[i] = Point(m.RawRowView(i))
	}
	return points

}
func colToPoints(m mat64.Dense) []Point {
	a := m.T()
	r, c := m.Dims()
	points := make([]Point, c)
	for i := range points {
		dst := make([]float64, r)
		points[i] = Point(mat64.Row(dst, i, a))
	}
	return points
}

type IG struct {
	I int
	G int
	V float64
	W float64
}
type IGs []IG

func (s IGs) Len() int {
	return len(s)
}
func (s IGs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/*
func (s IGs) Less(i, j int) bool {
	if s[i].G == s[j].G {
		return s[i].I < s[j].I
	}
	return s[i].G < s[j].G
}*/
func (s IGs) Less(i, j int) bool {
	if s[i].G == s[j].G {
		return s[i].W < s[j].W
	}
	return s[i].G < s[j].G
}

func weightPoint(x []float64) float64 {
	w := 0.0
	s := 0.0
	for i, v := range x {
		//log.Print(i, v)
		w += float64(i) * v
		s += v
	}
	//log.Println("weight", w/s)
	return w / s
}

func cos(x1 *mat64.Vector, x2 *mat64.Vector) float64 {
	n := x1.Len()
	p := mat64.NewVector(n, nil)
	p.MulElemVec(x1, x2)
	r := mat64.Sum(p) / (mat64.Norm(x1, 2) * mat64.Norm(x2, 2))
	if math.IsNaN(r) {
		r = 0.0
	}
	return r
}

/*
 * W r*k
 * H k*c
 * return cos r*c
 */
func cosMatrix(W *mat64.Dense, H *mat64.Dense) *mat64.Dense {
	r, _ := W.Dims()
	_, c := H.Dims()
	p := mat64.NewDense(r, c, nil)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			v1 := W.RowView(i)
			v2 := H.ColView(j)
			v := cos(v1, v2)
			p.Set(i, j, v)
		}
	}
	return p

}

/*
	mat r*c (cosMatrix output)

	return r*k  LDA scores (k groups )
           c*k  LDA scores
*/
func categoryLDAMat(mat *mat64.Dense, rowGroup []int, colGroup []int, k int) (*mat64.Dense, *mat64.Dense) {
	r := len(rowGroup)
	c := len(colGroup)
	rowarray := make([]float64, r*k)
	colarray := make([]float64, c*k)
	for i := 0; i < r; i++ {
		//log.Println(mat.RowView(i))
		lda := vectorCategoryLDAScore(mat.RowView(i), colGroup, k)
		for j := 0; j < k; j++ {
			rowarray[i*k+j] = lda[j]
		}
	}
	for i := 0; i < c; i++ {
		lda := vectorCategoryLDAScore(mat.ColView(i), rowGroup, k)
		for j := 0; j < k; j++ {
			colarray[i*k+j] = lda[j]
		}
	}
	rmat := mat64.NewDense(r, k, rowarray)
	cmat := mat64.NewDense(c, k, colarray)
	//log.Println(rmat)
	//log.Println(cmat)
	return rmat, cmat

}

func vectorCategoryLDAScore(v *mat64.Vector, cate []int, k int) []float64 {
	r := make([]float64, k)
	for i := 0; i < k; i++ {
		x0, x1 := splitCatVec(v, cate, i) //x[0] is not in this category, x[1] is in this category
		var d0 stats.Float64Data = x0
		var d1 stats.Float64Data = x1
		m0, err := stats.Mean(d0)
		checkErr(err)
		m1, err := stats.Mean(d1)
		checkErr(err)
		mu2 := (m1 - m0) * (m1 - m0)
		s0, err := stats.SampleVariance(d0)
		if err != nil {
			s0 = 0.0
		}
		s1, err := stats.SampleVariance(d1)
		if err != nil {
			s1 = 0.0
		}
		s := s1 + s0
		if s != 0.0 {
			r[i] = mu2 / s
		} else {
			if mu2 == 0.0 {
				r[i] = 0.0
			} else {
				r[i] = math.MaxFloat64
			}
		}
	}
	return r
}

/*
 *  lda = (u1-u2)^2/(s1+s2)
 */
func categoryLDAScore(v []float64, cate []int, k int) []float64 {
	r := make([]float64, k)
	for i := 0; i < k; i++ {
		x0, x1 := splitCat(v, cate, i) //x[0] is not in this category, x[1] is in this category
		var d0 stats.Float64Data = x0
		var d1 stats.Float64Data = x1
		m0, err := stats.Mean(d0)
		checkErr(err)
		m1, err := stats.Mean(d1)
		checkErr(err)
		mu2 := (m1 - m0) * (m1 - m0)
		s0, err := stats.SampleVariance(d0)
		if err != nil {
			s0 = 0.0
		}
		s1, err := stats.SampleVariance(d1)
		if err != nil {
			s1 = 0.0
		}
		s := s1 + s0
		if s != 0.0 {
			r[i] = mu2 / s
		} else {
			if mu2 == 0.0 {
				r[i] = 0.0
			} else {
				r[i] = math.MaxFloat64
			}
		}
	}
	return r
}

func splitCat(v []float64, cate []int, i0 int) ([]float64, []float64) {
	l := len(v)
	x0 := make([]float64, 0, l)
	x1 := make([]float64, 0, l)
	for i := 0; i < l; i++ {
		if cate[i] == i0 {
			x1 = append(x1, v[i])
		} else {
			x0 = append(x0, v[i])
		}
	}
	return x0, x1
}

func splitCatVec(v *mat64.Vector, cate []int, i0 int) ([]float64, []float64) {
	l := v.Len()
	x0 := make([]float64, 0, l)
	x1 := make([]float64, 0, l)
	//log.Println(l)
	for i := 0; i < l; i++ {
		if cate[i] == i0 {
			x1 = append(x1, v.At(i, 0))
		} else {
			x0 = append(x0, v.At(i, 0))
		}
	}
	return x0, x1
}

/*
 * get average for each category
 */
func categoryAverage(v []float64, cate []int, k int) []float64 {
	l := make([]float64, k)
	s := make([]float64, k)
	for i := 0; i < len(v); i++ {
		j := cate[i]
		l[j] += 1.0
		s[j] += v[i]
	}
	a := make([]float64, k)
	for i := 0; i < k; i++ {
		a[i] = s[i] / l[i]
	}
	return a
}

type IV struct {
	I int
	V float64
}
type IVs []IV

func (s IVs) Len() int {
	return len(s)
}
func (s IVs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IVs) Less(i, j int) bool {
	if s[i].V == s[j].V {
		return s[i].I < s[j].I
	}
	return s[i].V < s[j].V
}

func sortRank(v []float64) []int {
	arr := IVs(make([]IV, len(v)))
	for i := range v {
		arr[i] = IV{I: i, V: v[i]}
	}
	sort.Sort(arr)
	o := make([]int, len(v))
	for i := range arr {
		o[arr[i].I] = i
	}
	return o
}

func SortLabeledNMF(m *mat64.Dense, rowIds []string, colIds []string, k int) (*mat64.Dense, []string, []int, []string, []int, *mat64.Dense, *mat64.Dense, float64, *mat64.Dense, *mat64.Dense) {
	W, H, d, _ := NMF(m, k)
	rowPoints := rowToPoints(*W)
	colPoints := colToPoints(*H)
	r, c := m.Dims()

	//rowCenter, rowGroup , rowGroupNum, rowDistort := KMPP(rowPoints,rowGroupNum)
	_, rowGroup, rowGroupNum, _ := KMPP(rowPoints, k)
	_, colGroup, colGroupNum, _ := KMPP(colPoints, k)
	wW := make([]float64, r)

	wH := make([]float64, c)

	for i := 0; i < r; i++ {
		//log.Println(W.RowView(i).RawVector().Data)
		//log.Println("row view length", i, len(W.RowView(i).RawVector().Data))
		wW[i] = weightPoint(W.RowView(i).RawVector().Data)
	}

	for i := 0; i < c; i++ {
		a := make([]float64, k)
		v := H.ColView(i)
		for j := 0; j < k; j++ {
			a[j] = v.At(j, 0)
		}
		//log.Println(a)
		wH[i] = weightPoint(a)
	}
	//log.Println("H")
	//log.Println(H.Dims())

	kW := categoryAverage(wW, rowGroup, k)
	kH := categoryAverage(wH, colGroup, k)

	mapRowGroup := sortRank(kW)
	newRowGroupNum := make([]int, k)
	for i := 0; i < k; i++ {
		newRowGroupNum[mapRowGroup[i]] = rowGroupNum[i]
	}
	mapColGroup := sortRank(kH)
	newColGroupNum := make([]int, k)
	for i := 0; i < k; i++ {
		newColGroupNum[mapColGroup[i]] = colGroupNum[i]
	}
	//log.Println(kW, mapRowGroup)
	//log.Println(kH, mapColGroup)
	//TODO HERE
	//log.Println("kW", kW, rowGroupNum)
	//log.Println("kH", kH, colGroupNum)
	retv := mat64.NewDense(r, c, nil)
	rig := IGs(make([]IG, len(rowPoints)))
	for i := range rowPoints {
		rig[i] = IG{I: i, G: mapRowGroup[rowGroup[i]], V: kW[rowGroup[i]], W: wW[i]}
	}
	sort.Sort(rig)
	cig := IGs(make([]IG, len(colPoints)))
	for j := range colPoints {
		cig[j] = IG{I: j, G: mapColGroup[colGroup[j]], V: kH[colGroup[j]], W: wH[j]}
	}
	sort.Sort(cig)
	for i := range rig {
		i0 := rig[i].I
		for j := range cig {
			j0 := cig[j].I
			retv.Set(i, j, m.At(i0, j0))
		}
	}

	newW := mat64.NewDense(r, k, nil)

	newH := mat64.NewDense(k, c, nil)

	for i := range rig {
		for j := 0; j < k; j++ {
			newW.Set(i, j, W.At(rig[i].I, j))
		}
	}

	for i := range cig {
		for j := 0; j < k; j++ {
			newH.Set(j, i, H.At(j, cig[i].I))
		}
	}

	newColIds := make([]string, c)
	newColGroup := make([]int, c)
	newRowIds := make([]string, r)
	newRowGroup := make([]int, r)
	for i := range rig {
		newRowIds[i] = rowIds[rig[i].I]
		newRowGroup[i] = mapRowGroup[rowGroup[rig[i].I]]
	}
	for i := range cig {
		newColIds[i] = colIds[cig[i].I]
		newColGroup[i] = mapColGroup[colGroup[cig[i].I]]
	}
	//log.Println(newColGroup, newRowGroup)
	cosMat := cosMatrix(newW, newH)
	//log.Println(cosMat)
	a, b := categoryLDAMat(cosMat, newRowGroup, newColGroup, k)

	return retv, newRowIds, newRowGroup, newColIds, newColGroup, newW, newH, d, a, b
}
