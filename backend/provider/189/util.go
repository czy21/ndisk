package _189

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func hmacSha1(data string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func aesEncrypt(data, key []byte) string {
	block, _ := aes.NewCipher(key)
	if block == nil {
		return hex.EncodeToString([]byte{})
	}
	data = pkcs7Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}
	return hex.EncodeToString(decrypted)
}

func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func random(v string) string {
	reg := regexp.MustCompilePOSIX("[xy]")
	data := reg.ReplaceAllFunc([]byte(v), func(msg []byte) []byte {
		var i int64
		t := int64(16*rand.Float32()) | 0
		if msg[0] == 120 {
			i = t
		} else {
			i = 3&t | 8
		}
		return []byte(strconv.FormatInt(i, 16))
	})
	return string(data)
}

func rsaEncode(plainText []byte, publicKey string) string {
	publicKey = strings.Join([]string{"-----BEGIN PUBLIC KEY-----", publicKey, "-----END PUBLIC KEY-----"}, "\n")
	block, _ := pem.Decode([]byte(publicKey))
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubInterface.(*rsa.PublicKey)
	b, err := rsa.EncryptPKCS1v15(rand2.Reader, pub, plainText)
	if err != nil {
		log.Errorf("err: %s", err.Error())
	}
	return base64.StdEncoding.EncodeToString(b)
}
