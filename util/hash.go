package util

import (
	"crypto/md5" //nolint:gosec
	"crypto/sha256"
	"encoding/hex"
)

// ShaHash 返回 sha hash 字符串
func ShaHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}

// Md5 摘要，保证唯一性
func Md5(s string) string {
	sum := md5.Sum([]byte(s)) //nolint:gosec
	return hex.EncodeToString(sum[:])
}

func Md5FromBytes(s []byte) string {
	sum := md5.Sum(s) //nolint:gosec
	return hex.EncodeToString(sum[:])
}
