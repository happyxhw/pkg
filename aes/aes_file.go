package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

// EncryptFile encrypt file
func EncryptFile(in, out string, key []byte, fixedIV bool) error {
	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	var iv []byte
	if fixedIV {
		iv = key[:block.BlockSize()]
	} else {
		iv = make([]byte, block.BlockSize())
		_, err = io.ReadFull(rand.Reader, iv)
		if err != nil {
			return err
		}
	}

	outFile, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := inFile.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			_, _ = outFile.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	if !fixedIV {
		_, _ = outFile.Write(iv)
	}

	return nil
}

func DecryptFile(in, out string, key []byte, fixedIV bool) error {
	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	fi, err := inFile.Stat()
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	var iv []byte
	var msgLen int64
	if fixedIV {
		iv = key[:block.BlockSize()]
	} else {
		iv = make([]byte, block.BlockSize())
		msgLen = fi.Size() - int64(len(iv))
		_, err = inFile.ReadAt(iv, msgLen)
		if err != nil {
			return err
		}
	}

	outFile, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := inFile.Read(buf)
		if n > 0 {
			if !fixedIV {
				if n > int(msgLen) {
					n = int(msgLen)
				}
				msgLen -= int64(n)
			}
			stream.XORKeyStream(buf, buf[:n])
			_, _ = outFile.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
