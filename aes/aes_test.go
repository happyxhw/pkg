package aes

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptCBC(t *testing.T) {
	key, err := GenAesKey(32)
	require.NoError(t, err)

	source := "hello world"
	encrypted, err := EncryptCBC([]byte(source), key, nil)
	require.NoError(t, err)

	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println(encryptedB64)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	require.NoError(t, err)

	decrypted, err := DecryptCBC(encrypted, key, nil)
	require.NoError(t, err)
	fmt.Println(string(decrypted))
}

func TestEncryptGCM(t *testing.T) {
	key, err := GenAesKey(32)
	require.NoError(t, err)

	source := "hello world"
	encrypted, err := EncryptGCM([]byte(source), key)
	require.NoError(t, err)

	encryptedB64 := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println(encryptedB64)
	encrypted, err = base64.StdEncoding.DecodeString(encryptedB64)
	require.NoError(t, err)

	decrypted, err := DecryptGCM(encrypted, key)
	require.NoError(t, err)
	fmt.Println(string(decrypted))
}
