package aes

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenRsaKey(t *testing.T) {
	privatePath, publicPath := "./id_rsa", "./id_rsa.pub"
	err := GenRsaKey(4096, privatePath, publicPath, nil)
	require.NoError(t, err)
}

func TestRsaEncrypt(t *testing.T) {
	privatePath, publicPath := "./id_rsa", "./id_rsa.pub"
	err := GenRsaKey(4096, privatePath, publicPath, nil)
	require.NoError(t, err)

	privateKey, err := ReadPrivateKey(privatePath, nil)
	require.NoError(t, err)
	publicKey, err := ReadPublicKey(publicPath)
	require.NoError(t, err)

	source := "hello world"
	encrypted, err := RsaEncrypt([]byte(source), publicKey)
	require.NoError(t, err)
	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := RsaDecrypt(encrypted, privateKey)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decrypted))
}

func TestRsaEncryptWithPassword(t *testing.T) {
	passwd := []byte("123")
	privatePath, publicPath := "./id_rsa", "./id_rsa.pub"
	err := GenRsaKey(4096, privatePath, publicPath, passwd)
	require.NoError(t, err)

	privateKey, err := ReadPrivateKey(privatePath, passwd)
	require.NoError(t, err)
	publicKey, err := ReadPublicKey(publicPath)
	require.NoError(t, err)

	source := "hello world"
	encrypted, err := RsaEncrypt([]byte(source), publicKey)
	require.NoError(t, err)
	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := RsaDecrypt(encrypted, privateKey)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decrypted))
}
