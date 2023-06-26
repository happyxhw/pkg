package aes

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// GenRsaKey 生成密钥对
func GenRsaPrivateKey(bits int, privatePath string, pwd []byte) error {
	// 1. 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	// 2. MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derPrivateStream := x509.MarshalPKCS1PrivateKey(privateKey)

	// 3. Block代表PEM编码的结构, 对其进行设置
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateStream,
	}

	// 4. 创建文件
	privateFile, err := os.Create(privatePath)
	if err != nil {
		return err
	}
	defer func() { _ = privateFile.Close() }()

	if pwd != nil {
		data := pem.EncodeToMemory(block)
		encryptData, err2 := EncryptGCM(data, Gen256KeyFromPassword(pwd))
		if err2 != nil {
			return err2
		}
		err = os.WriteFile(privatePath, encryptData, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		err = pem.Encode(privateFile, block)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenRsaPublicKey(bits int, privatePath, publicPath string, pwd []byte) error {
	privateKey, err := ReadPrivateKey(privatePath, pwd)
	if err != nil {
		return err
	}

	// 1. 生成公钥文件
	publicKey := privateKey.PublicKey
	derPublicStream := x509.MarshalPKCS1PublicKey(&publicKey)

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPublicStream,
	}

	publicFile, err := os.Create(publicPath)
	if err != nil {
		return err
	}
	defer func() { _ = publicFile.Close() }()

	// 2. 编码公钥, 写入文件
	err = pem.Encode(publicFile, block)
	if err != nil {
		return err
	}
	return nil
}

func ReadPublicKey(keyPath string) (*rsa.PublicKey, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("invalid key")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func ReadPrivateKey(keyPath string, pwd []byte) (*rsa.PrivateKey, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	if pwd != nil {
		key, err = DecryptGCM(key, Gen256KeyFromPassword(pwd))
		if err != nil {
			return nil, err
		}
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("invalid key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func RsaEncrypt(src []byte, pk *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pk, src)
}

func RsaDecrypt(src []byte, pk *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, pk, src)
}
