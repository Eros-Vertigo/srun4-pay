package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

// Encrypt 加密给定的输入字符串
func Encrypt(data []byte, key string) (string, error) {
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 补全数据长度至8的倍数
	padding := des.BlockSize - (len(data) % des.BlockSize)
	paddedData := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)

	cipherText := make([]byte, len(paddedData))
	// 使用自定义的ECB加密器
	mode := NewECBEncrypter(block)
	mode.CryptBlocks(cipherText, paddedData)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 自定义 ECB 加密器
type ecbEncrypter struct {
	b         cipher.Block
	blockSize int
}

func NewECBEncrypter(block cipher.Block) cipher.BlockMode {
	return &ecbEncrypter{
		b:         block,
		blockSize: block.BlockSize(),
	}
}

func (e *ecbEncrypter) BlockSize() int {
	return e.blockSize
}

func (e *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("crypto/des: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/des: output smaller than input")
	}
	for i := 0; i < len(src); i += e.blockSize {
		e.b.Encrypt(dst[i:i+e.blockSize], src[i:i+e.blockSize])
	}
}
