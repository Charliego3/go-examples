package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kataras/golog"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/whimthen/temp/csvs"
	"github.com/whimthen/temp/times"
)

type Replacer func(i int, col string) string

var (
	replaceCmd      *cobra.Command
	csvfile         string
	timecols        []int
	decimalCols     []int
	fieldsPerRecord int
	lazyQuotes      bool

	replaceies = make(map[int]Replacer, 0)
)

func init() {
	replaceCmd = &cobra.Command{
		Use:   "replace",
		Short: "replace the csv file of the column",
		Run:   replace,
	}

	flags := replaceCmd.PersistentFlags()
	flags.StringVarP(&csvfile, "file", "f", "", "replace the csv file column value")
	flags.IntSliceVarP(&timecols, "timecols", "t", []int{}, "this is csv file change millsecond to parsedTime columns, like: 1,2,3")
	flags.IntSliceVarP(&decimalCols, "decimalCols", "d", []int{}, "the column of changed to decimal types")
	flags.BoolVarP(&lazyQuotes, "lazyQuotes", "l", true, "csv file reade lazy quote")
	flags.IntVarP(&fieldsPerRecord, "fieldsPerRecord", "r", 0, "csv filee read fieldsPerRecord")

	// entrust
	entrustCmd := &cobra.Command{
		Use: "entrust",
		Run: func(cmd *cobra.Command, args []string) {
			timecols = []int{11, 12}
			decimalCols = []int{1, 2, 3, 4, 5, 13}
			lazyQuotes = true
			replace(cmd, args)
		},
	}
	entrustCmd.PersistentFlags().StringVarP(&csvfile, "file", "f", "", "the entrust csv file")
	replaceCmd.AddCommand(entrustCmd)

	// trans_record
	transCmd := &cobra.Command{
		Use: "trans",
		Run: func(cmd *cobra.Command, args []string) {
			timecols = []int{9, 10}
			decimalCols = []int{1, 2, 3}
			lazyQuotes = true
			replace(cmd, args)
		},
	}
	transCmd.PersistentFlags().StringVarP(&csvfile, "file", "f", "", "the entrust csv file")
	replaceCmd.AddCommand(transCmd)

	// golog.SetLevel(golog.DebugLevel.String())
	golog.SetLevel(golog.InfoLevel.String())
}

func registerReplacer(i int, replacer Replacer) {
	replaceies[i] = replacer
}

func replace(cmd *cobra.Command, args []string) {
	golog.Infof("Will Replce TimeCols: %v, DecimalCols: %v", timecols, decimalCols)
	if len(timecols) <= 0 || len(decimalCols) <= 0 {
		golog.Warnf("No replace column of the csv file")
		return
	}
	for _, tc := range timecols {
		registerReplacer(tc, func(i int, col string) string {
			golog.Debugf("Column: [%d], Col: %s, ParsedTime: %s", i, col, times.Parse2I2S(col))
			ts := times.Parse2I2S(col)
			return ts
		})
		golog.Debugf("Registered column[%d] to time format", tc)
	}

	for _, dc := range decimalCols {
		registerReplacer(dc, func(i int, col string) string {
			f, err := decimal.NewFromString(col)
			if err != nil {
				golog.Errorf("Decimal format error: %v", err)
				return col
			}
			return f.String()
		})
	}

	if csvfile == "" {
		golog.Error("the csv file is not Specify!!!!")
		return
	}

	stat, err := os.Stat(csvfile)
	if err != nil || os.IsNotExist(err) {
		golog.Error(err)
		return
	}

	if stat.IsDir() {
		golog.Error("the csvfile is a dir, not a file")
		return
	}

	ext := filepath.Ext(csvfile)
	replaceFile := csvfile[:len(csvfile)-len(ext)] + "_replace" + ext
	_, err = os.Stat(replaceFile)
	if err == nil {
		os.Remove(replaceFile)
	}

	file, err := os.OpenFile(replaceFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		golog.Error("Can not create replace file")
		return
	}
	defer file.Close()

	csvs.Read(csvfile)
	csvs.Reader.FieldsPerRecord = fieldsPerRecord
	csvs.Reader.LazyQuotes = lazyQuotes

	// the Each func param i is line number from the csv file
	// the record is line data of each column
	_, err = csvs.Each(func(i int, record []string) {
		for j, col := range record {
			replacer := replaceies[j]
			if replacer != nil {
				record[j] = replacer(j, col)
			}
		}
		fmt.Fprintln(file, "\""+strings.Join(record, "\",\"")+"\"")
	})

	if err != nil {
		golog.Error(err)
		return
	}
	golog.Info("Replace csv file Successfuly!!!!")
}
