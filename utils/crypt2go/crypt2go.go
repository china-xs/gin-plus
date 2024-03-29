// Package crypt2go
// @file      : crypt2go.go
// @author    : china.gdxs@gmail.com
// @time      : 2024/2/3 16:45
// @Description:
package crypt2go

import (
	"crypto/aes"
	"github.com/china-xs/gin-plus/utils/crypt2go/ecb"
	"github.com/china-xs/gin-plus/utils/crypt2go/padding"
)

// AesEncrypt aes-加密
// @param plaintext
// @param key
// @return []byte
// @return error
func AesEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	plaintext, err = padder.Pad(plaintext) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		return nil, err
	}
	ct := make([]byte, len(plaintext))
	mode.CryptBlocks(ct, plaintext)
	return ct, nil
}

// AesDecrypt aes-解密
// @param plaintext
// @param key
// @return []byte
// @return error
func AesDecrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(plaintext))
	mode.CryptBlocks(pt, plaintext)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	return padder.Unpad(pt) // unpad plaintext after decryption
}
