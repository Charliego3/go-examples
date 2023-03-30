package main

import (
	"regexp"
	"strings"

	"github.com/charliego93/websocket"
)

type Config struct {
	Websockets []OrderConfig `yaml:"websockets"`
}

type OrderConfig struct {
	Name      string `yaml:"name"`
	ForURL    string `yaml:"for"`
	OwnURL    string `yaml:"own"`
	ApiURL    string `yaml:"api"`
	TrdURL    string `yaml:"trade"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Market    string `yaml:"market"`
}

type OkexReqMsg struct {
	*websocket.JsonMessage
	Op   string    `json:"op"`
	Args []OkexArg `json:"args"`
}

type OkexArg struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type OwnReqMsg struct {
	*websocket.JsonMessage
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

func NewOwnTickerReq(market string) *OwnReqMsg {
	return &OwnReqMsg{
		Event:   "addChannel",
		Channel: ClearMarket(market) + "_ticker",
	}
}

func NewOwnReq(channel string) *OwnReqMsg {
	return &OwnReqMsg{
		Event:   "addChannel",
		Channel: channel,
	}
}

func ClearMarket(market string) string {
	regex := regexp.MustCompile("[/_-]")
	market = regex.ReplaceAllString(market, "")
	return strings.ToLower(market)
}
