package main

import (
	_ "embed"
	"encoding/json"
	"github.com/kataras/golog"
)

var (
	//go:embed env.json
	content []byte
	Envs    map[string]string
)

func init() {
	Envs = make(map[string]string)
	if len(content) <= 0 {
		return
	}

	err := json.Unmarshal(content, &Envs)
	if err != nil {
		golog.Error("解析环境失败")
	}
}
