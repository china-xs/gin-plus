// Package utils
// @file      : p.jwt_test.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/12/20 21:15
// @Description:
package utils

import (
	"log"
	"testing"
)

func TestJWT_generateJWTSecret(t *testing.T) {
	var secret string
	var err error
	log.Printf("len:%v\n", len(`AAFlgk8SADB1VPoxvR9J_MRQto55IdNy2u8Mqj4J3e8CpXFxZNYWsIKgeeSMYLD1zdA4XQY7a0E`))
	if secret, err = generateJWTSecret(64); err != nil {
		t.Error(err)
	}
	log.Printf("secret:%s", secret)
}
