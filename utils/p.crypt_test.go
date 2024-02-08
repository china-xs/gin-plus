// Package utils
// @file      : p.crypt_test.go
// @author    : china.gdxs@gmail.com
// @time      : 2024/2/8 12:31
// @Description:
package utils

import "testing"

func TestAesEn2Str(t *testing.T) {
	// mysql 方案
	// SELECT TO_BASE64(AES_ENCRYPT('13429030111', '2R[<)NcD)^H4GDv.'));
	key := "2R[<)NcD)^H4GDv."
	obj := NewCryptByKey(key)
	var (
		str   string
		err   error
		phone string
	)
	phone = `13427848776`
	if str, err = obj.AesEn2Str(phone); err != nil {
		t.Error(err)
		return
	}
	if str, err = obj.AesDe2Str(str); err != nil {
		t.Error(err)
		return
	}

	if str != phone {
		t.Error(`解码错误`)
	}

}
