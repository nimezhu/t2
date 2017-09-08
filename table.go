package main

import (
	"math"
	//"bufio"
	"bytes"
	"encoding/csv"
	//"encoding/gob"
	//"errors"
	"fmt"
	"os"
	"strconv"
	//"strings"

	"github.com/gonum/matrix/mat64"
)

type Table struct {
	ColNames []string
	RowNames []string
	ColSize  int
	RowSize  int
	Mat      []float64
	FileName string
	Name     string
}

func NewTable(r int, c int, data []float64) *Table {
	colNames := make([]string, r)
	rowNames := make([]string, c)
	for i := range colNames {
		colNames[i] = "C" + strconv.Itoa(i)
	}
	for i := range rowNames {
		rowNames[i] = "R" + strconv.Itoa(i)
	}
	if data == nil {
		data = make([]float64, r*c)
		for i := range data {
			data[i] = 0.0
		}
	}
	fileName := "noname.tsv"
	name := "table"
	return &Table{colNames, rowNames, r, c, data, fileName, name}

}
func (t *Table) Dims() (int, int) {
	return t.RowSize, t.ColSize
}
func (t *Table) Cols() []string {
	return t.ColNames
}
func (t *Table) Rows() []string {
	return t.RowNames
}
func (t *Table) Dense() *mat64.Dense {
	return mat64.NewDense(t.RowSize, t.ColSize, t.Mat)
}
func (t *Table) String() string {
	return t.PrettyString(-1)
}
func (t *Table) TxtEncode() string {
	return t.PrettyString(2)
}

/*
func (t *Table) SaveGob(fileName string) error {
	bm, err0 := os.Create(fileName)
	if err0 != nil {
		return err0
	}
	enc := gob.NewEncoder(bm)
	err := enc.Encode(t)
	if err != nil {
		return err
	}
	bm.Close()
	return nil
}
func (t *Table) LoadGob(fileName string) error {
	bm, err0 := os.Open(fileName)
	defer bm.Close()
	if err0 != nil {
		return err0
	}
	dec := gob.NewDecoder(bm)
	err := dec.Decode(&t)
	if err != nil {
		return err
	}
	return nil
}
func (t *Table) Dims() (int, int) {
	return t.rowSize, t.colSize
}

/*
func (t *Table) SaveTsv(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(t.String())
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}
*/
func (t *Table) T() error {
	m := t.Dense().RawMatrix().Data
	t.ColNames, t.RowNames = t.RowNames, t.ColNames
	t.ColSize, t.RowSize = t.RowSize, t.ColSize
	t.FileName = t.FileName + "_transpose.tsv"
	t.Name = t.Name + "_transpose"
	data := make([]float64, t.RowSize*t.ColSize)
	for i := 0; i < t.ColSize; i++ {
		for j := 0; j < t.RowSize; j++ {
			data[j*t.ColSize+i] = m[i*t.RowSize+j]
		}
	}
	t.Mat = data
	return nil
}
func (t *Table) Info() string {
	r, c := t.Dims()
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Dims: %d * %d\n\nRowNames:\n\t", r, c))
	k := r
	s := ""
	if k > 20 {
		k = 20
		s = "..."
	}
	for i := 0; i < k; i++ {
		buffer.WriteString(t.RowNames[i] + ",")
	}
	buffer.WriteString(s)
	buffer.WriteString("\n\nColNames:\n\t")
	s = ""
	k = c
	if k > 20 {
		k = 20
		s = "..."
	}
	for j := 0; j < k; j++ {
		buffer.WriteString(t.ColNames[j] + ",")
	}
	buffer.WriteString(s)
	buffer.WriteString("\n\n")
	mat := t.Dense()
	max := mat64.Max(mat)
	min := mat64.Min(mat)
	buffer.WriteString(fmt.Sprintf("Domain: [%f , %f]", min, max))
	buffer.WriteString("\n\n")
	return buffer.String()
}

func (t *Table) Log(e float64, pseudo float64) error {
	a := make([]float64, len(t.Mat))
	root := math.Log(e)
	for i, v := range t.Mat {
		a[i] = math.Log(v+pseudo) / root
	}
	t.Mat = a
	t.Name += "|log"
	return nil
}
func (table *Table) LoadFile(f *os.File) error {
	r := csv.NewReader(f)
	r.Comma = '\t'
	table.FileName = f.Name()
	iter, err := r.ReadAll()
	if err != nil {
		return err
	}
	table.Name = iter[0][0]
	table.ColNames = iter[0][1:]
	table.ColSize = len(table.ColNames)
	table.RowSize = len(iter) - 1
	table.RowNames = make([]string, table.RowSize)
	table.Mat = make([]float64, table.ColSize*table.RowSize)
	for i := 1; i < len(iter); i++ {
		name, values := iter[i][0], iter[i][1:]
		for j := 0; j < len(values); j++ {
			table.Mat[(i-1)*table.ColSize+j], err = strconv.ParseFloat(values[j], 64)
			if err != nil {
				return err
			}
		}
		table.RowNames[i-1] = name
	}
	return err
}
func (table *Table) LoadTsv(file string) error {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}
	err = table.LoadFile(f)
	//table := new(Table)
	return err
}

func (t *Table) PrettyStringChosenRows(rows []int, f int) string {
	_, c := t.Dims()
	var buffer bytes.Buffer
	buffer.WriteString(t.Name + "_sub")
	for i0 := 0; i0 < c; i0++ {
		s := fmt.Sprintf("\t%s", t.ColNames[i0])
		buffer.WriteString(s)
	}

	buffer.WriteString("\n")
	format := "\t" + "%." + strconv.Itoa(f) + "f"
	if f == -1 {
		format = "\t%f"
	}
	m := t.Dense()
	for i := 0; i < len(rows); i++ {
		buffer.WriteString(fmt.Sprintf("%s", t.RowNames[rows[i]]))
		for j := 0; j < c; j++ {
			buffer.WriteString(fmt.Sprintf(format, m.At(rows[i], j)))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()

}
func (t *Table) PrettyString(f int) string {
	r, c := t.Dims()
	var buffer bytes.Buffer
	buffer.WriteString(t.Name)

	for i0 := 0; i0 < c; i0++ {
		s := fmt.Sprintf("\t%s", t.ColNames[i0])
		buffer.WriteString(s)
	}

	buffer.WriteString("\n")
	format := "\t" + "%." + strconv.Itoa(f) + "f"
	if f == -1 {
		format = "\t%f"
	}
	m := t.Dense()
	for i := 0; i < r; i++ {
		buffer.WriteString(fmt.Sprintf("%s", t.RowNames[i]))
		for j := 0; j < c; j++ {
			buffer.WriteString(fmt.Sprintf(format, m.At(i, j)))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
