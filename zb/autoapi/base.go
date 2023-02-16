package autoapi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/whimthen/temp/logger"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	DefaultAccount *Account
	client         = resty.New().SetTimeout(time.Minute).SetLogger(logger.StandardLogger()).SetProxy("http://172.16.100.150:23128")
)

type Config struct {
	Current  string    `yaml:"current"`
	Accounts []Account `yaml:"accounts"`
}

type Account struct {
	Account   string `yaml:"account"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	API       string `yaml:"api"`
	Trade     string `yaml:"trade"`
	KLine     string `yaml:"kline"`
	WSAPI     string `yaml:"wsapi"`
}

func request[T any](endpoint string, opts ...Option[*Values]) T {
	p := getOpts(opts...)
	p.Set("method", filepath.Base(endpoint))
	p.Set("accesskey", p.AccessKey)
	digestSign(p)
	p.Set("reqTime", fmt.Sprint(time.Now().Unix()*1000))
	if !strings.HasSuffix(p.URL, "/") {
		p.URL += "/"
	}
	t := new(T)
	request := p.URL + endpoint + "?" + p.Encode()
	resp, err := client.R().SetResult(t).Get(request)
	var cancel bool
	if err != nil {
		logger.Errorf("Request: %s\n\terror: %+v", request, err)
		cancel = true
	} else if resp.StatusCode() != 200 {
		logger.Infof("Request: %s\n\t Status: %s", request, resp.Status())
		cancel = true
	} else if !strings.HasPrefix(resp.Header().Get("Content-Type"), "application/json") {
		logger.Errorf("响应非JSON格式, Request: %s, Response: %s", request, resp.Body())
		cancel = true
	}
	if cancel && !p.continueErr {
		os.Exit(1)
	}
	return *t
}

func getOpts(opts ...Option[*Values]) *Values {
	p := &Values{Values: make(url.Values)}
	for _, opt := range opts {
		opt(p)
	}
	if p.Account == nil {
		p.Account = DefaultAccount
	}
	if p.URL == "" {
		if p.usingTrade {
			p.URL = p.Trade
		} else {
			p.URL = p.API
		}
	}
	return p
}

func loadConfig() {
	buf, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	config := Config{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	for _, acct := range config.Accounts {
		acct := acct
		if acct.Account == config.Current {
			DefaultAccount = &acct
			break
		}
	}
}
