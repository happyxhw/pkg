package aes

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncryptCBC(t *testing.T) {
	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	source := "hello world"
	encrypted, err := EncryptCBC([]byte(source), key, nil)
	if err != nil {
		t.Error(err)
	}
	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println(encryptedB64)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := DecryptCBC(encrypted, key, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decrypted))
}

func TestEncryptGCM(t *testing.T) {
	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	source := "hello world"
	encrypted, err := EncryptGCM([]byte(source), key)
	if err != nil {
		t.Error(err)
	}
	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println(encryptedB64)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := DecryptGCM(encrypted, key)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decrypted))
}
