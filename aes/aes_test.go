package aes

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	source := "hello world"
	encrypted, err := Encrypt([]byte(source), key, nil)
	if err != nil {
		t.Error(err)
	}
	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println(encryptedB64)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := Decrypt(encrypted, key, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decrypted))
}
