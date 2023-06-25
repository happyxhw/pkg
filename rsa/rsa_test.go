package rsa

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenKey(t *testing.T) {
	privatePath, publicPath := "./id_rsa", "./id_rsa.pub"
	err := GenKey(4096, privatePath, publicPath)
	require.NoError(t, err)
}

func TestEncrypt(t *testing.T) {
	privatePath, publicPath := "./id_rsa", "./id_rsa.pub"
	err := GenKey(4096, privatePath, publicPath)
	require.NoError(t, err)

	source := "hello world"
	encrypted, err := Encrypt([]byte(source), publicPath)
	require.NoError(t, err)
	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := Decrypt(encrypted, privatePath)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decrypted))
}
