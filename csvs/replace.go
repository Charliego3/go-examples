package csvs

import "github.com/sirupsen/logrus"

type ReplacerFunc func(d string) string

var rrs = make(map[int]ReplacerFunc)

func RegisterFunc(f ReplacerFunc, col ...int) {
	for _, i := range col {
		rrs[i] = f
	}
}

func Replace(j int) [][]string {
	records, err := Each(func(i int, record []string) {
		if i >= j {
			for ri, r := range record {
				if f, ok := rrs[ri]; ok {
					record[ri] = f(r)
				}
			}
		}
	})

	if err != nil {
		logrus.Error("Replace error", err)
		return nil
	}

	return records
}
