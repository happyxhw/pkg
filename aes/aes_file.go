package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const bufLen = 32 * 1024

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

	buf := make([]byte, bufLen)
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

	buf := make([]byte, bufLen)
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

func EncryptFileWithRSA(in, out string, keyKey []byte, pk *rsa.PublicKey) error {
	// 随机生成 key 和 iv
	key, err := GenAesKey(32)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	iv, err := GenAesKey(block.BlockSize())
	if err != nil {
		return err
	}

	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	tmp := make([]byte, 0, len(key)+len(iv))
	tmp = append(tmp, key...)
	tmp = append(tmp, iv...)
	// aes 加密 key iv
	tmp, err = EncryptGCM(tmp, keyKey)
	if err != nil {
		return err
	}
	// rsa 加密 key iv
	encTmp, err := RsaEncrypt(tmp, pk)
	if err != nil {
		return err
	}
	var head []byte
	n := len(encTmp)
	head = append(head, byte(n), byte(n>>8))
	head = append(head, encTmp...)
	_, err = outFile.Write(head)
	if err != nil {
		return err
	}

	buf := make([]byte, bufLen)
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

	return nil
}

func DecryptFileWithRSA(in, out string, keyKey []byte, pk *rsa.PrivateKey) error {
	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	tmp := make([]byte, bufLen)
	_, err = io.ReadFull(inFile, tmp[:2])
	if err != nil {
		return err
	}
	n := int(tmp[0]) | int(tmp[1])<<8
	if n > bufLen {
		return errors.New("len(rsa) out of index")
	}
	_, err = io.ReadFull(inFile, tmp[:n])
	if err != nil {
		return err
	}
	// rsa 解密
	keyIV, err := RsaDecrypt(tmp[:n], pk)
	if err != nil {
		return err
	}
	// aes 解密
	keyIV, err = DecryptGCM(keyIV, keyKey)
	if err != nil {
		return err
	}
	key, iv := keyIV[:32], keyIV[32:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	buf := make([]byte, bufLen)
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
	return nil
}

func EncryptFileAndPathWithRSA(in, out string, keyKey []byte, pk *rsa.PublicKey) error {
	// 随机生成 key 和 iv
	key, err := GenAesKey(32)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	iv, err := GenAesKey(block.BlockSize())
	if err != nil {
		return err
	}

	inFile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	tmp := make([]byte, 0, len(key)+len(iv))
	tmp = append(tmp, key...)
	tmp = append(tmp, iv...)
	// aes 加密 key iv
	tmp, err = EncryptGCM(tmp, keyKey)
	if err != nil {
		return err
	}
	// rsa 加密 key iv
	encTmp, err := RsaEncrypt(tmp, pk)
	if err != nil {
		return err
	}
	var head []byte
	n := len(encTmp)
	head = append(head, byte(n), byte(n>>8))
	head = append(head, encTmp...)
	_, err = outFile.Write(head)
	if err != nil {
		return err
	}

	var pathHead []byte
	encPath, err := EncryptGCM([]byte(in), keyKey)
	if err != nil {
		return err
	}
	n = len(encPath)
	pathHead = append(pathHead, byte(n), byte(n>>8))
	pathHead = append(pathHead, encPath...)
	_, err = outFile.Write(pathHead)
	if err != nil {
		return err
	}

	buf := make([]byte, bufLen)
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

	return nil
}

func DecryptFileAndPathWithRSA(in, out string, keyKey []byte, pk *rsa.PrivateKey) (string, error) {
	inFile, err := os.Open(in)
	if err != nil {
		return "", err
	}
	defer inFile.Close()

	tmp := make([]byte, bufLen)
	_, err = io.ReadFull(inFile, tmp[:2])
	if err != nil {
		return "", err
	}
	n := int(tmp[0]) | int(tmp[1])<<8
	if n > bufLen {
		return "", errors.New("len(rsa) out of index")
	}
	_, err = io.ReadFull(inFile, tmp[:n])
	if err != nil {
		return "", err
	}
	// rsa 解密
	keyIV, err := RsaDecrypt(tmp[:n], pk)
	if err != nil {
		return "", err
	}
	// aes 解密
	keyIV, err = DecryptGCM(keyIV, keyKey)
	if err != nil {
		return "", err
	}
	key, iv := keyIV[:32], keyIV[32:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	tmp = make([]byte, 2)
	_, err = io.ReadFull(inFile, tmp[:2])
	if err != nil {
		return "", err
	}
	n = int(tmp[0]) | int(tmp[1])<<8
	pathEnc := make([]byte, n)
	_, err = io.ReadFull(inFile, pathEnc[:n])
	if err != nil {
		return "", err
	}
	pathDec, err := DecryptGCM(pathEnc, keyKey)
	if err != nil {
		return "", err
	}
	if out == "" {
		out = string(pathDec)
		dir := filepath.Dir(out)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	outFile, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	buf := make([]byte, bufLen)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := inFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		stream.XORKeyStream(buf, buf[:n])
		if _, wErr := outFile.Write(buf[:n]); wErr != nil {
			return "", wErr
		}
	}
	return string(pathDec), nil
}
