package autoapi

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func digestSign(p *Values) {
	sign := HmacMD5(p.Encode(), p.SecretKey)
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
