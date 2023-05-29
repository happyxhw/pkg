package fsutil

import (
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"io"
	"os"
)

func GetFileMD5(in string) (string, error) {
	inFile, err := os.Open(in)
	if err != nil {
		return "", err
	}
	defer inFile.Close()

	md5h := md5.New() //nolint:gosec
	_, err = io.Copy(md5h, inFile)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(md5h.Sum(nil)), nil
}
