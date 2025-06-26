package aes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptFile(t *testing.T) {
	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	err := EncryptFile("./aes_file.go", "./out.txt", key, true)
	require.NoError(t, err)

	err = DecryptFile("./out.txt", "./in_2.txt", key, true)
	require.NoError(t, err)
}

func TestEncryptFileWithRSA(t *testing.T) {
	privatePath, publicPath := "./id_rsa", "./id_rsa.pub"
	require.NoError(t, GenRsaPrivateKey(4096, privatePath, nil))
	require.NoError(t, GenRsaPublicKey(4096, privatePath, publicPath, nil))

	privateKey, err := ReadPrivateKey(privatePath, nil)
	require.NoError(t, err)
	publicKey, err := ReadPublicKey(publicPath)
	require.NoError(t, err)

	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	err = EncryptFileWithRSA("./aes_file.go", "./out.txt", key, publicKey)
	require.NoError(t, err)

	err = DecryptFileWithRSA("./out.txt", "./in_2.txt", key, privateKey)
	require.NoError(t, err)
}

func TestEncryptFileAndPathWithRSA(t *testing.T) {
	privatePath, publicPath := "./id_rsa", "./id_rsa.pub"
	require.NoError(t, GenRsaPrivateKey(4096, privatePath, nil))
	require.NoError(t, GenRsaPublicKey(4096, privatePath, publicPath, nil))

	privateKey, err := ReadPrivateKey(privatePath, nil)
	require.NoError(t, err)
	publicKey, err := ReadPublicKey(publicPath)
	require.NoError(t, err)

	key := []byte("cQfTjWnZr4u7x!A%D*G-KaPdRgUkXp2s")
	err = EncryptFileAndPathWithRSA("./in.txt", "./out.txt", key, publicKey)
	require.NoError(t, err)

	path, err := DecryptFileAndPathWithRSA("./out.txt", "", key, privateKey)
	require.NoError(t, err)
	require.Equal(t, path, "./in.txt")
}
