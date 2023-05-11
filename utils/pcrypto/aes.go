// AES 加解密
package pcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/perpower/goframe/funcs/normal"
)

var Aes = gaes{}

// 定义AES 结构体
type gaes struct{}

// Encrypt AES加密 初始向量16字节空 PKCS7 CBC
// origData: string 待加密串
// key:密钥 string 16/24/32
// iv: 向量 []byte
// 返回:加密后 string
func (c *gaes) Encrypt(origData, key string, iv []byte) (string, error) {
	_data := normal.String2Bytes(origData)
	_key := normal.String2Bytes(key)

	if !validKey(_key) {
		panic("秘钥长度错误")
	}

	block, err := aes.NewCipher(_key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	//补全码
	_data = pkcs7Padding(_data, blockSize)

	dst := make([]byte, len(_data))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(dst, _data)

	return base64.StdEncoding.EncodeToString(dst), nil
}

// Decrypt AES解密 初始向量16字节空 PKCS7 CBC
// 入参:origData string 加密字符串
// key:密钥 string 16/24/32位
// iv: 向量 []byte
// 返回:解密后 string
func (c *gaes) Decrypt(origData, key string, iv []byte) (string, error) {
	_data, _ := base64.StdEncoding.DecodeString(origData)
	_key := normal.String2Bytes(key)

	if !validKey(_key) {
		panic("秘钥长度错误")
	}

	block, err := aes.NewCipher(_key)
	if err != nil {
		return "", err
	}

	// blockSize := block.BlockSize()
	dst := make([]byte, len(_data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, _data)

	//去除补全码
	dst = pkcs7UnPadding(dst)

	return normal.Bytes2String(dst), nil
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)

	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 秘钥长度验证
func validKey(key []byte) bool {
	k := len(key)
	switch k {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}
