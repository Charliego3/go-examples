package main

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/shopspring/decimal"
	"github.com/whimthen/kits/logger"
	"strings"
)

func main() {
	file, err := excelize.OpenFile("/Users/nzlong/Downloads/用户ID：499116现货账单.xlsx")
	if err != nil {
		logger.Fatal("Can't open the xls file", err)
	}

	sheetList := file.GetSheetList()
	for _, sn := range sheetList {
		logger.Debug("Sheet Name: %s", sn)
		rows, err := file.Rows(sn)
		if err != nil {
			logger.Fatal("Get rows from %s error: %+v", sn, err)
		}

		err = rows.Error()
		if err != nil {
			logger.Fatal("Rows error: %+v", err)
		}

		var prec, prea *decimal.Decimal
		var pred bool
		var precol []string

		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				logger.Fatal("Rows column error: %+v", err)
			}

			changeQc := columns[2]
			var change, amount string

			if strings.Contains(changeQc, "<br>") {
				for _, s := range strings.Split(changeQc, "<br>") {
					if strings.HasSuffix(s, "QC") {
						changeQc = s
					}
				}
			}

			if !strings.HasSuffix(changeQc, "QC") {
				continue
			}
			changeQc = strings.TrimSpace(strings.ReplaceAll(changeQc, "QC", ""))

			vs := strings.Split(changeQc, "=")
			if len(vs) != 2 {
				logger.Warn("Columns: %+v", columns)
				logger.Fatal("ChangeQC value length != 2")
			}

			change = strings.TrimSpace(vs[0])
			amount = strings.TrimSpace(vs[1])

			isDeduct := false
			if strings.HasPrefix(change, "-") {
				isDeduct = true
				change = change[1:]
			}
			if strings.HasPrefix(change, "+") {
				change = change[1:]
			}

			if strings.HasPrefix(amount, "+") {
				amount = amount[1:]
			}

			//logger.Warn("IsDeduct: %v, Change: %s, Amount: %s", isDeduct, change, amount)

			c := toDecimal(change)
			a := toDecimal(amount)

			var isNotEqual bool
			if prec != nil && prea != nil {
				var r decimal.Decimal
				if pred {
					r = prec.Add(*prea)
				} else {
					r = prea.Sub(*prec)
				}
				isNotEqual = !r.Equal(*a)

				if isNotEqual {
					logger.Error("计算结果不一致, \n上一行: %+v, \n当前行: %+v", precol, columns)
				}
			}

			prec = c
			prea = a
			pred = isDeduct
			precol = columns
			continue
		}
	}
}

func toDecimal(v string) *decimal.Decimal {
	d, err := decimal.NewFromString(v)
	if err != nil {
		logger.Fatal("Format value to decimal error: %v", err)
	}
	return &d
}
