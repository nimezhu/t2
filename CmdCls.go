package main

import (
	"fmt"
	"log"

	"github.com/nimezhu/fire"
	"github.com/urfave/cli"
)

func CmdCls(cli *cli.Context) {
	k := cli.Int("k")
	t, _ := loadTsv(cli)
	r, c := t.Dims()
	mat, rowIds, rowGroup, colIds, colGroup, _, _, _, _, _ := fire.SortLabeledNMF(t.Dense(), t.Rows(), t.Cols(), k)
	log.Println(colGroup)
	log.Println(rowGroup)
	newMat := Table{colIds, rowIds, c, r, mat.RawMatrix().Data, t.FileName + "_cls", t.Name + "_cls"}
	fmt.Print(newMat.PrettyString(-1))
}
