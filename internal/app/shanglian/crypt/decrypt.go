package crypt

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
)

// Decrypt 解密给定的加密字符串
func Decrypt(encrypted, key string) ([]byte, error) {
	// 解码加密数据
	ciphertext, _ := base64.StdEncoding.DecodeString(encrypted)

	// 创建一个 DES 解密器
	block, _ := des.NewCipher([]byte(key))

	// 解密数据
	decrypted := make([]byte, len(ciphertext))
	blockMode := NewECBDecrypter(block)
	blockMode.CryptBlocks(decrypted, ciphertext)

	// 去除 PKCS5 填充
	decrypted = pkcs5Unpad(decrypted)

	fmt.Println("解密结果:", string(decrypted))

	return decrypted, nil
}

func pkcs5Unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

// 自定义 ECB 解密器
type ecbDecrypter struct {
	b         cipher.Block
	blockSize int
}

func NewECBDecrypter(block cipher.Block) cipher.BlockMode {
	return &ecbDecrypter{
		b:         block,
		blockSize: block.BlockSize(),
	}
}

func (d *ecbDecrypter) BlockSize() int {
	return d.blockSize
}

func (d *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%d.blockSize != 0 {
		panic("crypto/des: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/des: output smaller than input")
	}
	for i := 0; i < len(src); i += d.blockSize {
		d.b.Decrypt(dst[i:i+d.blockSize], src[i:i+d.blockSize])
	}
}
