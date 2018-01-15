package main

import (
	"fmt"
	"log"
	"strings"

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
	s := newMat.PrettyString(-1) //TODO
	lines := strings.Split(s, "\n")
	header := strings.Split(lines[0], "\t")
	fmt.Printf("%s\t%s\t%s\n", header[0], "rowGroup", strings.Join(header[1:], "\t"))
	for i, line := range lines[1:] {
		f := strings.Split(line, "\t")
		if len(line) == 0 {
			continue
		}
		fmt.Printf("%s\t%d\t%s\n", f[0], rowGroup[i], strings.Join(f[1:], "\t"))
	}
	//TODO
}
