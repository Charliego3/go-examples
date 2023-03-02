package autoapi

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strings"
)

func encodeDigestSign(p *Values) {
	sign := HmacMD5(p.Encode(), p.SecretKey)
	p.Set("sign", sign)
}

func digestSign(p *Values) {
	var params strings.Builder
	keys := make([]string, 0, len(p.Values))
	for k := range p.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := p.Values[k]
		for _, v := range vs {
			if params.Len() > 0 {
				params.WriteByte('&')
			}
			params.WriteString(k)
			params.WriteByte('=')
			params.WriteString(v)
		}
	}
	sign := HmacMD5(params.String(), p.SecretKey)
	p.Set("sign", sign)
}

// HmacMD5 MD5
func HmacMD5(message, secretKey string) (hmacSign string) {
	h := hmac.New(md5.New, []byte(digest(secretKey)))
	h.Write([]byte(message))
	hmacSign = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA1 加密
func digest(secretKey string) (digest string) {
	hash := sha1.New()
	hash.Write([]byte(secretKey))
	digest = hex.EncodeToString(hash.Sum(nil))
	return
}
