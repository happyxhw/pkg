package aes

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGen256Key(t *testing.T) {
	key, err := Gen256Key()
	require.NoError(t, err)
	encryptedB64 := base64.StdEncoding.EncodeToString(key)
	fmt.Printf("debug-x: %s\n", encryptedB64)
}
