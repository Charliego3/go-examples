package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	cookieHeader = `JSESSIONID=2a82002e-24ac-4e78-a60f-5556d0319a84; org.springframework.web.servlet.i18n.CookieLocaleResolver.LOCALE=zh_CN; CASTGC=TGT-3219-11DaYdtctvBaNoJN1A2RfoqFFYmL4e4pbKlJQ3CEcYQKM2RNHu-sso14; casusername=%E7%89%9B%E8%B6%85%E9%BE%99; tgw_l7_route=4929821a7bf17b68601c3b2550afc97a; ariafontScale=true; ariawapChangeViewPort=false`

	headers = map[string]string{
		//"Accept-Encoding":           `gzip, deflate, br`,
		"Accept-Language":           `zh-CN,zh;q=0.9`,
		"Cache-Control":             `max-age=0`,
		"Connection":                `keep-alive`,
		"Content-Length":            `363`,
		"Content-Type":              `application/x-www-form-urlencoded`,
		"Cookie":                    `JSESSIONID=2a82002e-24ac-4e78-a60f-5556d0319a84; org.springframework.web.servlet.i18n.CookieLocaleResolver.LOCALE=zh_CN; CASTGC=TGT-3219-11DaYdtctvBaNoJN1A2RfoqFFYmL4e4pbKlJQ3CEcYQKM2RNHu-sso14; casusername=%E7%89%9B%E8%B6%85%E9%BE%99; tgw_l7_route=4929821a7bf17b68601c3b2550afc97a; ariafontScale=true; ariawapChangeViewPort=false`,
		"Host":                      `msjw.ga.sz.gov.cn`,
		"Origin":                    `https://msjw.ga.sz.gov.cn`,
		"Referer":                   `https://msjw.ga.sz.gov.cn/crj/crjmsjw/wsyy/ajax/scheduleSzjm`,
		"sec-ch-ua":                 `"Not_A Brand";v="99", "Google Chrome";v="109", "Chromium";v="109"`,
		"sec-ch-ua-mobile":          `?0`,
		"sec-ch-ua-platform":        `macOS`,
		"Sec-Fetch-Dest":            `document`,
		"Sec-Fetch-Mode":            `navigate`,
		"Sec-Fetch-Site":            `same-origin`,
		"Sec-Fetch-User":            `?1`,
		"Upgrade-Insecure-Requests": `1`,
		"User-Agent":                `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36`,
	}

	client = resty.New().
		SetBaseURL("https://msjw.ga.sz.gov.cn/crj/crjmsjw/").
		SetHeaders(headers).
		SetTimeout(time.Minute)
)

type Reponse[T any] struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type Reserve struct {
	Sfrc     bool `json:"sfrc"`
	DataList []struct {
		Date    string `json:"date"`
		Col     int    `json:"col"`
		Num     int    `json:"num"`
		Context string `json:"context"`
		Xq      int    `json:"xq"`
		Time    string `json:"time"`
		Row     int    `json:"row"`
	} `json:"dataList"`
	Rclx      interface{} `json:"rclx"`
	RowNumMap struct {
		RowNum int `json:"rowNum"`
	} `json:"rowNumMap"`
}

func getReserve() {

}

func main() {
	var r Reponse[Reserve]
	client.R().SetFormData(map[string]string{
		"yywdbh":    "440309000000",
		"blyw":      "-1",
		"startDate": "2023-02-07",
		"type":      "1",
		"days":      "7",
		"lang":      "CH",
		"ywlx":      "szjm",
	}).SetResult(&r).Post("wsyydata/getScheduleData")
	log.Println(r)

	// var rs string
	// var err any
	// client.R().SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8").
	// 	SetError(&err).
	// 	SetResult(&rs).Get("wsyy/ajax/szjm.html")
	// log.Println(rs)
	// log.Println(err)
	login()
}

func login() {
	request, err := http.NewRequest(http.MethodGet, "https://msjw.ga.sz.gov.cn/crj/crjmsjw/wsyy/ajax/szjm.html", nil)
	if err != nil {
		log.Fatal(err)
	}

	for key, val := range headers {
		request.Header.Add(key, val)
	}
	request.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`)

	client := &http.Client{Timeout: time.Minute}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost,
		"https://msjw.ga.sz.gov.cn/szga_yhtx_cas/login?loginType=1&locale=zh_CN&service=https%3A%2F%2Fmsjw.ga.sz.gov.cn%2Fcrj%2Fcrjmsjw%2Fwsyy%2Fajax%2Fszjm.html", nil)
	for key, val := range headers {
		req.Header.Add(key, val)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	_ = res

	if response.StatusCode == 302 {
		log.Panicln("跳转")
	}
}
