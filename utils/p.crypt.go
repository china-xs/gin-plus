// Package utils
// @file      : p.crypt.go
// @author    : china.gdxs@gmail.com
// @time      : 2024/2/3 16:49
// @Description:
package utils

import (
	"encoding/base64"
	"github.com/china-xs/gin-plus/utils/crypt2go"
	"github.com/spf13/viper"
)

type Crypt struct {
	key []byte
}

func NewCrypt(viper *viper.Viper) *Crypt {
	keyAny := viper.Get(`cryptKey`)
	var (
		key string
		ok  bool
	)
	if key, ok = keyAny.(string); !ok || key == `` {
		panic(`cryptKey is null `)
	}
	return &Crypt{
		key: []byte(key),
	}
}

func NewCryptByKey(key string) *Crypt {
	return &Crypt{
		key: []byte(key),
	}
}

// AesEn2Str aes str 加密
func (this *Crypt) AesEn2Str(str string) (result string, err error) {
	var buf []byte
	if buf, err = crypt2go.AesDecrypt([]byte(str), this.key); err != nil {
		return
	}
	result = base64.StdEncoding.EncodeToString(buf)
	return
}

// AesDe2Str 手机号码解析
func (this *Crypt) AesDe2Str(str string) (result string, err error) {
	var buf []byte
	if buf, err = base64.StdEncoding.DecodeString(str); err != nil {
		return
	}
	if buf, err = crypt2go.AesDecrypt(buf, this.key); err != nil {
		return
	}
	result = string(buf)
	return
}
