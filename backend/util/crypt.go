package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Encode(plainText string) string {
	h := md5.New()
	h.Write([]byte(plainText))
	return hex.EncodeToString(h.Sum(nil))
}

func GetMd5Bytes(p []byte) []byte {
	h := md5.New()
	h.Write(p)
	return h.Sum(nil)
}
