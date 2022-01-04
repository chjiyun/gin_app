// 加解密工具包
package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func Encrypt(src string, key []byte) string {
	data := []byte(src)
	block, err := des.NewCipher(key[:8])
	if err != nil {
		panic(err)
	}
	data = PKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[8:16])
	out := make([]byte, len(data))
	blockMode.CryptBlocks(out, data)
	return fmt.Sprintf("%x", out)
}

// des解密, key长度必须大于16位
func Decrypt(src string, key []byte) string {
	data, _ := hex.DecodeString(src)
	block, err := des.NewCipher(key[:8])
	if err != nil {
		panic(err)
	}
	// 实例化解密模式(参数为密码对象和密钥)
	blockMode := cipher.NewCBCDecrypter(block, key[8:16])
	str := make([]byte, len(data))
	blockMode.CryptBlocks(str, data)
	str = PKCS5UnPadding(str)
	return string(str)
}

//明文补码算法
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//明文减码算法
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Sha1 散列函数之sha1
func Sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	// 最终的散列值的字符切片，参数可以用来对现有的字符切片追加额外的字节切片
	out := h.Sum(nil)
	return fmt.Sprintf("%x", out)
}
