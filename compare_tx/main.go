package main

import (
	"bufio"
	"github.com/whimthen/kits/logger"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("/Users/nzlong/Downloads/0917-bitbank-cacel.txt")
	if err != nil {
		logger.Fatal("file no found.", err)
	}

	txIds := make(map[string]int)
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if len(line) > 0 {
			txId := strings.TrimSpace(string(line))
			t, ok := txIds[txId]
			if ok {
				txIds[txId] = t + 1
			} else {
				txIds[txId] = 1
			}
		}
	}

	defer file.Close()

	logger.Debug("%+v", txIds)
}
