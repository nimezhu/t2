package main

import (
	"fmt"

	"github.com/montanaflynn/stats"
	"github.com/nimezhu/fire"
	. "github.com/nimezhu/ice"
	"github.com/urfave/cli"
)

func CmdColStats(c *cli.Context) error {
	tsv, err := loadTsv(c)
	mat := tsv.Dense()
	rowNum, colNum := mat.Dims()
	fmt.Printf("label\tmin\t1st_quantile\tmedian\t3rd_quantile\tmax\tmean\tgini\n")
	for i := 0; i < colNum; i++ {
		data := make([]float64, rowNum)
		for j := 0; j < rowNum; j++ {
			data[j] = mat.At(j, i)
		}
		min, _ := stats.Min(data)
		q, _ := stats.Quartile(data)
		max, _ := stats.Max(data)
		mean, _ := stats.Mean(data)
		gini := fire.Gini(data)

		fmt.Printf("%s\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n", tsv.ColNames[i], min, q.Q1, q.Q2, q.Q3, max, mean, gini)
	}
	return err
}

/*
type Table struct {
	ColNames []string
	RowNames []string
	ColSize  int
	RowSize  int
	Mat      []float64
	FileName string
	Name     string
}
*/
func CmdColCorr(c *cli.Context) error {
	tsv, err := loadTsv(c)
	mat := tsv.Dense()
	rowNum, colNum := mat.Dims()
	cols := make([][]float64, colNum)
	cor := make([]float64, colNum*colNum)
	for i := 0; i < colNum; i++ {
		cols[i] = make([]float64, rowNum)
		for j := 0; j < rowNum; j++ {
			cols[i][j] = mat.At(j, i)
		}
	}
	for i := 0; i < colNum; i++ {
		cor[i*colNum+i] = 1.0
		for j := i + 1; j < colNum; j++ {
			c, _ := stats.Correlation(cols[i], cols[j])
			cor[i*colNum+j] = c
			cor[j*colNum+i] = c
		}
	}
	corTable := Table{tsv.ColNames, tsv.ColNames, tsv.ColSize, tsv.ColSize, cor, tsv.FileName + "_col_corr", tsv.Name + "_col_corr"}
	fmt.Print(corTable.PrettyString(-1))
	return err
}
