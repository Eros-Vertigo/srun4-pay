package crypt

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/pkcs12"
)

// parsePKCS12 从PKCS12数据中解析出私钥和证书
func parsePKCS12(p12Data []byte, password string) (*rsa.PrivateKey, *x509.Certificate, error) {
	// 解码PKCS12数据
	blocks, err := pkcs12ToPEM(p12Data, password)
	if err != nil {
		return nil, nil, err
	}

	var cert *x509.Certificate
	var privateKey *rsa.PrivateKey

	// 遍历PEM块
	for _, block := range blocks {
		// 解析私钥
		if block.Type == "PRIVATE KEY" {
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, err
			}
			privateKey = key
		}

		// 解析证书
		if block.Type == "CERTIFICATE" {
			certs, err := x509.ParseCertificates(block.Bytes)
			if err != nil {
				return nil, nil, err
			}
			if len(certs) > 0 {
				cert = certs[0]
			}
		}
	}

	if privateKey == nil || cert == nil {
		return nil, nil, fmt.Errorf("无法从PKCS12中提取私钥或证书")
	}

	return privateKey, cert, nil
}

// pkcs12ToPEM 将PKCS12数据转换为PEM块
func pkcs12ToPEM(data []byte, password string) ([]*pem.Block, error) {
	blocks, err := pkcs12.ToPEM(data, password)
	if err != nil {
		return nil, err
	}

	var pemBlocks []*pem.Block
	for _, block := range blocks {
		pemBlock := &pem.Block{
			Type:  block.Type,
			Bytes: block.Bytes,
		}
		pemBlocks = append(pemBlocks, pemBlock)
	}

	return pemBlocks, nil
}

// getMd5 计算字符串的MD5哈希值
func getMd5(strDes string) string {
	bytes := []byte(strDes)
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:])
}
