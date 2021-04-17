package csvs

import (
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"os"
)

var Reader *csv.Reader

func Read(file string) {
	f, err := os.Open(file)
	if err != nil {
		logrus.Error("打开文件失败", err)
		return
	}
	Reader = csv.NewReader(f)
}

func All() ([][]string, error) {
	return Reader.ReadAll()
}

func Each(f func(i int, record []string)) ([][]string, error) {
	records, err := Reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		f(i, record)
	}
	return records, nil
}
