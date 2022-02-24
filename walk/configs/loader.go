package configs

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
)

type config struct {
	UUID          string `json:"uuid,omitempty"`
	AppKey        string `json:"appKey,omitempty"`
	AppSecret     string `json:"appSecret,omitempty"`
	VersionCode   string `json:"versionCode,omitempty"`
	VersionName   string `json:"versionName,omitempty"`
	Timestap      string `json:"timestap,omitempty"`
	BrandIdentify string `json:"brandIdentify,omitempty"`
	BgSource      string `json:"bgSource,omitempty"`
	PartType      string `json:"partType,omitempty"`
	CatAppId      string `json:"catAppId,omitempty"`
}

var Config = config{}

func init() {
	bytes, err := ioutil.ReadFile("meituan/config.json")
	if err != nil {
		log.Fatalln("读取配置文件失败", err)
	}
	err = jsoniter.Unmarshal(bytes, &Config)
	if err != nil {
		log.Fatalln("反序列化配置文件失败")
	}
}
