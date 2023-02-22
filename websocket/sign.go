package websocket

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/whimthen/temp/logger"
	"gopkg.in/errgo.v2/errors"
)

func AesSha1Prng(key []byte, length int) ([]byte, error) {
	bs := Sha1(Sha1(key))
	ml := len(bs)
	rl := length / 8
	if rl > ml {
		return nil, errors.New("AesSha1Prng invalid length")
	}
	return bs[:rl], nil
}

func generateKey(key []byte) []byte {
	genK := make([]byte, 16)
	copy(genK, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genK[j] ^= key[i]
		}
	}
	return genK
}

func AesEncryptECB(src []byte, key []byte) (encrypted []byte) {
	key, _ = AesSha1Prng(key, 128)
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

func AesDecryptECB(encrypted, key []byte) (decrypted []byte) {
	defer func() {
		if err := recover(); err != nil {
			decrypted = nil
			logger.Errorf("AesDecryptECB[%s] panic: %+v", encrypted, err)
		}
	}()

	key, err := AesSha1Prng(key, 128)
	if err != nil {
		logger.Errorf("Encrypted: %s, Key: %s, AesSha1Prng error: %+v", encrypted, key, err)
		return nil
	}

	cipher, err := aes.NewCipher(generateKey(key))
	if err != nil {
		logger.Errorf("Encrypted: %s, Key: %s, NewCipher error: %+v", encrypted, key, err)
		return nil
	}

	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

func MD5(src []byte) []byte {
	hash := md5.New()
	hash.Write(src)
	return hash.Sum(nil)
}

func HexMD5(src []byte) string {
	return hex.EncodeToString(MD5(src))
}

func HmacSHA256(src, secret []byte) string {
	hash := hmac.New(sha256.New, secret)
	hash.Write(src)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func HmacMD5(src, secret []byte) string {
	h := hmac.New(md5.New, secret)
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(src []byte) []byte {
	hash := sha1.New()
	hash.Write(src)
	return hash.Sum(nil)
}

func HexSha1(src []byte) (digest []byte) {
	hash := sha1.New()
	hash.Write(src)
	src = hash.Sum(nil)
	digest = make([]byte, hex.EncodedLen(len(src)))
	_ = hex.Encode(digest, src)
	return
}
