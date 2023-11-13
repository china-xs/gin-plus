// Package utils
// @file      : p.viper.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/9 16:11
// @Description:
package utils

import (
	"github.com/spf13/viper"
)

// NewViper 初始化viper
func NewViper(path string) (v *viper.Viper, err error) {
	v.AddConfigPath(".")
	v.SetConfigFile(string(path))
	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, err
}
