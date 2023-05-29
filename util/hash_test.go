package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShaHash(t *testing.T) {
	shaHash := ShaHash("hello, world")

	assert.NotEmpty(t, shaHash)
}

func TestMd5(t *testing.T) {
	m := Md5("hello, world")
	assert.NotEmpty(t, m)
}
