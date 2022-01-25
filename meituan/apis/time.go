package apis

import (
	jsoniter "github.com/json-iterator/go"
	"time"
)

const timeURL = "http://apimobile.meituan.com/group/v1/timestamp/milliseconds"

var useLocal bool

type currentMS struct {
	CurrentMS int64 `json:"currentMs"`
}

func NowMillis() int64 {
	return time.Now().UnixMilli()
}

func Milliseconds() int64 {
	if useLocal {
		return NowMillis()
	}

	resp, err := client.R().SetHeaders(map[string]string{
		"Accept":     "application/json",
		"user-agent": "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)",
		"Host":       "apimobile.meituan.com",
		"Connection": "keep-alive",
	}).Get(timeURL)
	if err != nil {
		return NowMillis()
	}

	var ms currentMS
	err = jsoniter.Unmarshal(resp.Body(), &ms)
	if err != nil {
		return NowMillis()
	}
	millis := NowMillis()
	if millis+1000*15 > ms.CurrentMS && millis-1000*15 < ms.CurrentMS {
		useLocal = true
	}
	return ms.CurrentMS
}
