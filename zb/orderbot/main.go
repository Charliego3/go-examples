package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

var config = Config{}

func loadConfig() {
	bs, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal("读取配置文件失败", "err", err)
	}
	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		log.Fatal("配置解析失败", "err", err)
	}
}

func main() {
	// formatter := "15:04:05"
	// log.SetTimeFormat(formatter)
	log.SetLevel(log.DebugLevel)
	loadConfig()
	if len(config.Websockets) == 0 {
		log.Fatal("未配置websocket")
	}

	ctx, cancel := context.WithCancel(context.Background())
	for _, c := range config.Websockets {
		opts := log.Options{
			ReportTimestamp: true,
			ReportCaller:    true,
			TimeFormat:      time.DateTime,
			Level:           log.DebugLevel,
			Prefix:          c.Name,
		}
		go (&Orderer{ctx: ctx, opts: opts}).Start(c)
	}

	stoper := make(chan os.Signal, 3)
	signal.Notify(stoper, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-stoper
	cancel()
	time.Sleep(time.Second * 2)
	log.Info("程序已停止")
}
