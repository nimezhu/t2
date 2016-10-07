package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	//"time"

	. "github.com/nimezhu/ice"
	"github.com/urfave/cli"
)

const (
	VERSION = "0.0.1"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = VERSION
	app.Name = "table tools"
	app.Usage = "handle tsv files [ first column, first row are labels, the rest are float64 number]"
	//app.EnableBashCompletion = true
	// global level flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show more output",
		},
	}

	// Commands
	app.Commands = []cli.Command{
		{
			Name:   "binrows",
			Usage:  "merge rows into bin / change different resolution",
			Action: CmdBinRows,
		},
		{
			Name:   "info",
			Usage:  "report summary of table",
			Action: CmdInfo,
		},
		{
			Name:   "view",
			Usage:  "view table [precision]",
			Action: CmdView,
		},
		{
			Name:  "log",
			Usage: "log table [pseudo count]",
			Flags: []cli.Flag{
				cli.Float64Flag{
					Name:  "root,r",
					Value: math.E,
					Usage: "log root",
				},
				cli.Float64Flag{
					Name:  "pseudo,p",
					Value: 1.0,
					Usage: "pseudo count",
				},
			},
			Action: CmdLog,
		},
		{
			Name:   "trans",
			Usage:  "trans file.tsv",
			Action: CmdT,
		},
		{
			Name:   "gini",
			Usage:  "gini file.tsv",
			Action: CmdColGini,
		},
	}
	app.Run(os.Args)
}

func dN(n int, size int) int {
	if !(n%size == 0) {
		return n/size + 1

	} else {
		return n / size
	}
}

/*
	ColNames []string
	RowNames []string
	ColSize  int
	RowSize  int
	Mat      []float64
	FileName string
	Name     string
*/
func binRows(t *Table, size int) *Table {
	r, c := t.Dims()
	mat := t.Dense()
	newR := dN(r, size)
	arr := make([]float64, newR*c)
	newRowNames := make([]string, newR)
	for i := 0; i < newR; i++ {
		i1 := i * size
		newRowNames[i] = t.RowNames[i1]
		l1 := i1 + size
		if l1 > r {
			l1 = r
		}
		for j := 0; j < c; j++ {
			index := i*c + j
			num := 0.0
			s := 0.0
			for i0 := i1; i0 < l1; i0++ {
				s += mat.At(i0, j)
				num += 1.0
			}
			arr[index] = s / num
		}
	}
	newTable := Table{t.ColNames, newRowNames, c, newR, arr, t.FileName + "_rowbin" + strconv.Itoa(size), t.Name + "_rowbin" + strconv.Itoa(size)}
	return &newTable
}
func CmdBinRows(c *cli.Context) error {
	if c.NArg() == 0 {
		log.Fatal("Usage: binrow file.tsv [binsize, default:2]")
	}
	fn := c.Args().Get(0)
	binsize := 2
	if c.NArg() > 2 {
		bin, err := strconv.Atoi(c.Args().Get(2))
		checkErr(err)
		binsize = bin
	}
	tsv := Table{}
	err := tsv.LoadTsv(fn)
	checkErr(err)
	newTsv := binRows(&tsv, binsize)
	fmt.Print(newTsv.PrettyString(-1))
	return err
}

func CmdInfo(c *cli.Context) error {
	if c.NArg() == 0 {
		log.Fatal("Usage: info file.tsv")
	}
	fn := c.Args().Get(0)
	tsv := Table{}
	err := tsv.LoadTsv(fn)
	checkErr(err)
	fmt.Print(tsv.Info())
	return err
}

func CmdView(c *cli.Context) error {
	if c.NArg() == 0 {
		log.Fatal("Usage: view file.tsv [precision, default:-1]")
	}
	p := -1
	if c.NArg() > 1 {
		t, err := strconv.Atoi(c.Args().Get(1))
		if err != nil {
			p = -1
		} else {
			p = t
		}
	}
	fn := c.Args().Get(0)
	tsv := Table{}
	err := tsv.LoadTsv(fn)
	checkErr(err)
	fmt.Print(tsv.PrettyString(p))
	return err
}

func CmdLog(c *cli.Context) error {
	if c.NArg() == 0 {
		log.Fatal("Usage Help : log -h")
	}
	root := c.Float64("root")
	p := c.Float64("pseudo")
	fn := c.Args().Get(0)
	tsv := Table{}
	err := tsv.LoadTsv(fn)
	checkErr(err)
	tsv.Log(root, p)
	fmt.Print(tsv.PrettyString(-1))
	return err
}

func CmdT(c *cli.Context) error {
	if c.NArg() == 0 {
		log.Fatal("Usage Help : log -h")
	}
	fn := c.Args().Get(0)
	tsv := Table{}
	err := tsv.LoadTsv(fn)
	checkErr(err)
	tsv.T()
	fmt.Print(tsv.PrettyString(-1))
	return err
}

func CmdColGini(c *cli.Context) error {
	fn := c.Args().Get(0)
	tsv := Table{}
	err := tsv.LoadTsv(fn)
	checkErr(err)
	mat := tsv.Dense()
	rowNum, colNum := mat.Dims()
	fmt.Printf("col\tlabel\tgini\n")
	for i := 0; i < colNum; i++ {
		data := make([]float64, rowNum)
		for j := 0; j < rowNum; j++ {
			data[j] = mat.At(j, i)
		}
		fmt.Printf("%d\t%s\t%f\n", i+1, tsv.ColNames[i], Gini(data))
	}
	return nil

}
