package aes

import (
	"crypto/rand"
	"io"
)

// GenKey 生成密钥对
func GenAesKey(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Gen256KeyFromPassword 从密码生成 key
func Gen256KeyFromPassword(password []byte) []byte {
	key := make([]byte, 32)
	if len(password) > 32 {
		copy(key, password[:32])
	} else {
		copy(key, password)
	}

	return key
}
