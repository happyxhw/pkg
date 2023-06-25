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
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		stream.XORKeyStream(buf, buf[:n])
		if _, wErr := outFile.Write(buf[:n]); wErr != nil {
			return wErr
		}
	}
	if !fixedIV {
		if _, wErr := outFile.Write(iv); wErr != nil {
			return wErr
		}
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
	fileLen := fi.Size()
	if fixedIV {
		iv = key[:block.BlockSize()]
	} else {
		iv = make([]byte, block.BlockSize())
		fileLen -= int64(len(iv))
		_, err = inFile.ReadAt(iv, fileLen)
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
	remainingLen := fileLen
	for {
		if remainingLen <= 0 {
			break
		}
		n, err := inFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !fixedIV {
			if n > int(remainingLen) {
				n = int(remainingLen)
			}
			remainingLen -= int64(n)
		}
		stream.XORKeyStream(buf, buf[:n])
		if _, wErr := outFile.Write(buf[:n]); wErr != nil {
			return wErr
		}
	}
	return nil
}
