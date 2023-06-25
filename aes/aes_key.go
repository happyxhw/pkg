package aes

import "crypto/rand"

// GenKey 生成密钥对
func Gen256Key() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
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
