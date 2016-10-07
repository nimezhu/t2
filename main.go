package main

import (
	"fmt"
	//"log"
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
			Name:  "binrows",
			Usage: "merge rows into bin / change different resolution",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "binsize,b",
					Value: 2,
					Usage: "bin window size",
				},
			},
			Action: CmdBinRows,
		},
		{
			Name:   "info",
			Usage:  "report summary of table",
			Action: CmdInfo,
		},
		{
			Name:  "view",
			Usage: "view table",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "precision,p",
					Value: -1,
					Usage: "precision",
				},
			},
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
		{
			Name:   "stats",
			Usage:  "stats for column",
			Action: CmdColStats,
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
	binsize := c.Int("binsize")
	tsv, err := loadTsv(c)
	newTsv := binRows(tsv, binsize)
	fmt.Print(newTsv.PrettyString(-1))
	return err
}

func CmdInfo(c *cli.Context) error {
	tsv, err := loadTsv(c)
	fmt.Print(tsv.Info())
	return err
}

func loadTsv(c *cli.Context) (*Table, error) {
	tsv := Table{}
	if c.NArg() == 0 {
		err := tsv.LoadFile(os.Stdin)
		return &tsv, err
	} else {
		fn := c.Args().Get(0)
		err := tsv.LoadTsv(fn)
		checkErr(err)
		return &tsv, err
	}
}
func CmdView(c *cli.Context) error {
	tsv, err := loadTsv(c)
	p := c.Int("precision")
	checkErr(err)
	fmt.Print(tsv.PrettyString(p))
	return err
}

func CmdLog(c *cli.Context) error {
	tsv, err := loadTsv(c)
	root := c.Float64("root")
	p := c.Float64("pseudo")
	tsv.Log(root, p)
	fmt.Print(tsv.PrettyString(-1))
	return err
}

func CmdT(c *cli.Context) error {
	tsv, err := loadTsv(c)
	checkErr(err)
	tsv.T()
	fmt.Print(tsv.PrettyString(-1))
	return nil
}

func CmdColGini(c *cli.Context) error {
	tsv, err := loadTsv(c)
	mat := tsv.Dense()
	rowNum, colNum := mat.Dims()
	fmt.Printf("label\tgini\n")
	for i := 0; i < colNum; i++ {
		data := make([]float64, rowNum)
		for j := 0; j < rowNum; j++ {
			data[j] = mat.At(j, i)
		}
		fmt.Printf("%s\t%f\n", tsv.ColNames[i], Gini(data))
	}
	return err

}
