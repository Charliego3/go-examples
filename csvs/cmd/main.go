package main

import (
	"bytes"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whimthen/temp/csvs"
	"github.com/whimthen/temp/times"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cmd := &cobra.Command{
		Use:   "csvs",
		Short: "csvs is changed csv file of columns",
		Run:   runner,
	}

	cmd.AddCommand(replaceCmd)
	err := cmd.Execute()
	if err != nil {
		logrus.Error(err)
	}
}

func runner(amd *cobra.Command, args []string) {
	filename := "/Users/nzlong/Downloads/sqlresult_5650327.csv"
	csvs.Read(filename)
	csvs.Reader.FieldsPerRecord = 0
	csvs.Reader.LazyQuotes = true
	csvs.RegisterFunc(func(d string) string {
		return times.Parse2I2S(d)
	}, 11, 12)
	csvs.RegisterFunc(func(d string) string {
		f, err := decimal.NewFromString(d)
		if err != nil {
			logrus.Error("decimal 转化失败", err)
		}
		return f.String()
	}, 13)
	replaced := csvs.Replace(1)

	buffer := bytes.Buffer{}
	for _, s := range replaced {
		buffer.WriteString("\"")
		buffer.WriteString(strings.Join(s, "\",\""))
		buffer.WriteString("\"\n")
	}

	ext := filepath.Ext(filename)
	if err := os.WriteFile(filename[:len(filename)-len(ext)]+"_111111"+ext, buffer.Bytes(), os.ModePerm); err != nil {
		logrus.Error(err)
	}
}
