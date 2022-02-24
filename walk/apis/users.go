package apis

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/whimthen/temp/meituan/configs"
	"log"
	"strconv"
)

const (
	// "https://epassport.meituan.com/bizapi/loginv5?appkey=app-pos&bg_source=7&reqtime=1642756506&sign=d08dd73257830b7ed68aa9d770ed110f756260d6&utm_campaign=uisdk1.0&utm_medium=android&utm_term=3.20.600&uuid=6064b84d5dee4adfa9d89418ae1d91c0a164236619454453772"
	loginURL = "https://epassport.meituan.com/bizapi/loginv5?appkey=app-pos&bg_source=7&reqtime=%d&sign=%s&utm_campaign=uisdk1.0&utm_medium=android&utm_term=3.20.600&uuid=%s"
)

var (
	loginHeaders = map[string]string{
		"Accept-Encoding": "gzip",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Host":            "epassport.meituan.com",
		"Connection":      "Keep-Alive",
		"User-Agent":      "okhttp/2.7.5",
		"Accept-Content":  "application/json",
	}
)

type LoginResp struct {
	Data struct {
		Bizacctid        int    `json:"bizacctid"`
		Login            string `json:"login"`
		PartType         int    `json:"part_type"`
		PartKey          string `json:"part_key"`
		LoginSensitive   int    `json:"loginSensitive"`
		NameSensitive    int    `json:"nameSensitive"`
		ContactSensitive int    `json:"contactSensitive"`
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpireIn         int    `json:"expire_in"`
		RefreshIn        int    `json:"refresh_in"`
	} `json:"data"`
}

func LoginV5(username, password, partKey string, remember bool) (resp LoginResp, err error) {
	rememberPwd := 1
	if !remember {
		rememberPwd = 0
	}

	reqTime := Milliseconds() / 1000
	// password=yy12347890&remember_password=1&part_type=1&part_key=6519763&login=18929387993&fingerprint=fingerprint
	form := map[string]string{
		"password":          password,
		"remember_password": strconv.Itoa(rememberPwd),
		"part_type":         "1",
		"part_key":          partKey,
		"login":             username,
		"fingerprint":       "fingerprint",
	}
	sign := Sign(reqTime)

	parsedLoginURL := fmt.Sprintf(loginURL, reqTime, sign, configs.Config.UUID)
	response, err := client.R().SetHeaders(loginHeaders).
		SetFormData(form).Post(parsedLoginURL)
	if err != nil {
		return
	}

	body := response.Body()
	log.Println("LoginV5 resp:", string(body))
	err = jsoniter.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	return
}
