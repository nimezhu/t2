package main

import (
	"fmt"

	"github.com/montanaflynn/stats"
	//. "github.com/nimezhu/ice"
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
		gini := Gini(data)

		fmt.Printf("%s\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n", tsv.ColNames[i], min, q.Q1, q.Q2, q.Q3, max, mean, gini)
	}
	return err
}
