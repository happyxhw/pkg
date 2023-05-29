package util

import (
	"crypto/rand"
)

/*
Refer to go-nanoid project:
https://github.com/matoous/go-nanoid
*/

const defaultSize = 32

// defaultAlphabet is the alphabet used for ID characters by default.
var defaultAlphabet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NanoID(length int) string {
	if length <= 0 {
		length = defaultSize
	}
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	id := make([]rune, length)
	for i := 0; i < length; i++ {
		id[i] = defaultAlphabet[bytes[i]&61]
	}
	return string(id[:length])
}
