package aes

import "testing"

func TestEncryptFile(t *testing.T) {
	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	err := EncryptFile("./aes_file.go", "./out.txt", key, true)
	if err != nil {
		t.Error(err)
	}
}

func TestDecryptFile(t *testing.T) {
	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	err := DecryptFile("./out.txt", "./in_2.txt", key, true)
	if err != nil {
		t.Error(err)
	}
}
