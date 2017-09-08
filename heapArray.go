package main

import (
	"container/heap"
	//"errors"
	//"log"
	//	"strconv"
	"sort"

	"github.com/gonum/matrix/mat64"
)

type Ele struct {
	I int
	V float64
}

type EleArr []Ele

func (h EleArr) Len() int           { return len(h) }
func (h EleArr) Less(i, j int) bool { return h[i].V > h[j].V }
func (h EleArr) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *EleArr) Push(x interface{}) {
	*h = append(*h, x.(Ele))
}

func (h *EleArr) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func matColToArr(m *mat64.Dense, j int, r int) (*EleArr, error) {
	arr := make(EleArr, r)
	for i := 0; i < r; i++ {
		arr[i] = Ele{I: i, V: m.At(i, j)}
	}
	return &arr, nil
}
func TblTopK(t *Table, k int, s int, e int) ([]int, error) {
	mat := t.Dense()
	r, _ := mat.Dims()
	set := make(map[int]bool)
	for col := s; col < e; col++ {
		h, _ := matColToArr(mat, col, r)
		i := k
		heap.Init(h)
		for i > 0 {
			e, _ := heap.Pop(h).(Ele)
			i -= 1
			set[e.I] = true

			/* mark e.I */
		}
	}
	var keys []int
	for k := range set {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys, nil
}
