package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/dlgs"
	"github.com/kataras/golog"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	golog.SetLevel("debug")

	targetFile, b, err := dlgs.File("Select html xls", "", false)
	println(targetFile, b, err)
	if !b || err != nil {
		golog.Error(err)
	}

	//targetFile := "/Users/nzlong/Downloads/result_2021_6_18下午4_42_30.xls"
	file, err := os.OpenFile(targetFile, os.O_RDONLY, 0644)
	if err != nil {
		golog.Fatal(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		golog.Fatal(err)
	}

	buffer := bytes.Buffer{}
	doc.Find("table tr").Each(func(i int, r *goquery.Selection) {
		td := r.Find("td")
		rows := make([]string, td.Children().Size())
		td.Each(func(i int, d *goquery.Selection) {
			rows = append(rows, fmt.Sprintf("\"%s\"", strings.ReplaceAll(d.Text(), "\"", "\"\"")))
		})
		buffer.WriteString(strings.Join(rows, ","))
		buffer.WriteByte('\n')
	})

	ioutil.WriteFile(targetFile+".csv", buffer.Bytes(), os.ModePerm)
}
