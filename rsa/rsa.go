package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// GenKey 生成密钥对
func GenKey(bits int, privatePath, publicPath string) error {
	// 1. 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	// 2. MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derPrivateStream := x509.MarshalPKCS1PrivateKey(privateKey)

	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateStream,
	}

	// 4. 创建文件
	privateFile, err := os.Create(privatePath)
	defer func() { _ = privateFile.Close() }()

	if err != nil {
		return err
	}

	// 5. 使用pem编码, 并将数据写入文件中
	err = pem.Encode(privateFile, &block)
	if err != nil {
		return err
	}

	// 1. 生成公钥文件
	publicKey := privateKey.PublicKey
	derPublicStream := x509.MarshalPKCS1PublicKey(&publicKey)

	block = pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPublicStream,
	}

	publicFile, err := os.Create(publicPath)
	defer func() { _ = publicFile.Close() }()

	if err != nil {
		return err
	}

	// 2. 编码公钥, 写入文件
	err = pem.Encode(publicFile, &block)
	if err != nil {
		return err
	}
	return nil
}

func Encrypt(src []byte, keyPath string) ([]byte, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("invalid key")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 公钥加密
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
}

func Decrypt(src []byte, keyPath string) ([]byte, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// 从数据中解析出pem块
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("invalid key")
	}

	// 解析出一个der编码的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 私钥解密
	result, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	if err != nil {
		return nil, err
	}
	return result, nil
}
