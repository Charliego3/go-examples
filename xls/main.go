package main

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/kataras/golog"
	"strings"
)

func main() {
	golog.SetLevel("debug")

	targetFile := "/Users/nzlong/Downloads/122050.xls"
	csvFile := targetFile + ".csv"
	golog.Info(csvFile)

	if xlFile, err := xls.Open(targetFile, "utf-8"); err == nil {
		fmt.Println(xlFile.Author)
		//第一个sheet
		sheet := xlFile.GetSheet(0)
		builder := strings.Builder{}
		if sheet.MaxRow != 0 {
			//temp := make([][]string, sheet.MaxRow)
			for i := 0; i < int(sheet.MaxRow); i++ {
				row := sheet.Row(i)
				data := make([]string, 0)
				if row.LastCol() > 0 {
					for j := 0; j < row.LastCol(); j++ {
						col := row.Col(j)
						data = append(data, col)
					}
					//temp[i] = data
					builder.WriteString(strings.Join(data, ","))
					builder.WriteByte('\n')
				}
			}
			//res = append(res, temp...)
		}

		golog.Info(builder.String())
	} else {
		golog.Errorf("open_err: %v", err)
	}
}
