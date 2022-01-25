package apis

import (
	"bytes"
	"crypto/sha1"
	"github.com/go-resty/resty/v2"
	"github.com/whimthen/temp/meituan/configs"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var client = resty.New()

func Sign(reqTime int64, queries ...string) string {
	signMap := url.Values{
		"appkey":       []string{configs.Config.AppKey},
		"bg_source":    []string{configs.Config.BgSource},
		"reqtime":      []string{strconv.FormatInt(reqTime, 10)},
		"utm_campaign": []string{"uisdk1.0"},
		"utm_medium":   []string{"android"},
		"utm_term":     []string{configs.Config.VersionName},
		"uuid":         []string{configs.Config.UUID},
	}

	if len(queries) > 0 {
		qns := strings.SplitN(queries[0], "=", -1)
		for i := 0; i < len(qns); i = i + 2 {
			signMap[qns[i]] = []string{qns[i+1]}
		}
	}

	var keys []string
	for key := range signMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	buff := bytes.Buffer{}
	buff.WriteString(configs.Config.AppSecret)
	for _, key := range keys {
		val := signMap[key][0]
		if val == "" || val == "0" {
			continue
		}

		buff.WriteString(key + val)
	}

	hash := sha1.New()
	hash.Write(buff.Bytes())
	sums := hash.Sum(nil)
	return offset(sums)
}

func offset(arr []byte) string {
	builder := strings.Builder{}
	for _, b := range arr {
		b2 := (b >> 4) & 15
		i := 0
		for {
			if b2 < 0 || b2 > 9 {
				builder.WriteByte((b2 - 10) + 97)
			} else {
				builder.WriteByte(b2 + 48)
			}
			b2 = b & 15
			i2 := i + 1
			if i >= 1 {
				break
			}
			i = i2
		}
	}
	return strings.ToLower(builder.String())
}
