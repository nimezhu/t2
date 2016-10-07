package main

import (
	"fmt"
	"log"

	. "github.com/nimezhu/ice"
	"github.com/urfave/cli"
)

func CmdCls(cli *cli.Context) {
	k := cli.Int("k")
	t, _ := loadTsv(cli)
	r, c := t.Dims()

	//mat, rowIds, rowGroup, colIds, colGroup, W, H, _, rowLDA, colLDA := SortLabeledNMF(t.Dense(), t.Rows(), t.Cols(), k)

	mat, rowIds, rowGroup, colIds, colGroup, _, _, _, _, _ := SortLabeledNMF(t.Dense(), t.Rows(), t.Cols(), k)
	//colGroup := Table{[]string{"colGroup"}, colIds, 1, c, float64arr(colGroup), "colGroup", "colGroup"}

	//memory.Add(&Table{[]string{"rowGroup"}, rowIds, 1, r, float64arr(rowGroup), "rowGroup", "rowGroup"}, "rowGroup")
	//memory.Add(&Table{colIds, generateNames(k, "E"), c, k, H.RawMatrix().Data, "H", "H"}, "H")
	//memory.Add(&Table{generateNames(k, "E"), rowIds, k, r, W.RawMatrix().Data, "W", "W"}, "W")
	log.Println(colGroup)
	log.Println(rowGroup)

	newMat := Table{colIds, rowIds, c, r, mat.RawMatrix().Data, t.FileName + "_cls", t.Name + "_cls"}
	fmt.Print(newMat.PrettyString(-1))
	//memory.Add(&Table{generateNames(k, "Group"), rowIds, k, r, rowLDA.RawMatrix().Data, "rowLDA", "rowLDA"}, "rowLDA")
	//memory.Add(&Table{colIds, generateNames(k, "Group"), k, c, colLDA.RawMatrix().Data, "colLDA", "colLDA"}, "colLDA")
}
