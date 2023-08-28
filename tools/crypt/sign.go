package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

// Sign 使用私钥对字符串进行签名
func Sign(strDes, certPassword, pfx string) (string, error) {
	data := getMd5(strDes)

	// 读取私钥文件
	p12Data, err := os.ReadFile(pfx)
	if err != nil {
		return "", err
	}

	// 解析PKCS12数据
	privateKey, _, err := parsePKCS12(p12Data, certPassword)
	if err != nil {
		return "", err
	}
	hash := sha1.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	// 使用RSA私钥进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed)
	if err != nil {
		return "", err
	}

	// 对签名进行Base64编码
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	return encodedSignature, nil
}

// Verify 使用公钥对签名进行验证
func Verify(strDes, signMsg, certFile string) (bool, error) {
	data := getMd5(strDes)

	// 解码Base64编码的签名
	decodedSignMsg, err := base64.StdEncoding.DecodeString(signMsg)
	if err != nil {
		return false, err
	}

	// 读取公钥文件
	publicKeyData, err := os.ReadFile(certFile)
	if err != nil {
		return false, err
	}

	// 解析公钥证书
	cert, err := parsePublicKeyCertificate(publicKeyData)
	if err != nil {
		return false, err
	}

	// 使用RSA公钥进行验签
	hash := sha1.New() // 使用SHA-1哈希算法
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA1, hashed, decodedSignMsg)
	if err != nil {
		return false, err
	}

	return true, nil
}

// parsePublicKeyCertificate parses the provided public key certificate data.
func parsePublicKeyCertificate(publicKeyData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(publicKeyData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode public key certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
