// Package utils
// @file      : providers.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/9 16:08
// @Description: utils
package utils

import "github.com/google/wire"

const (
	_tranceId = `traceId`
	_spanId   = `spanId`
)

var ProviderSet = wire.NewSet(
	NewViper, // config *viper.Viper
	NewLog,   // log *zapLogger
	NewGorm,
	NewGormLog,
	NewRedis,
	NewJwtAuth,
	NewJWTOps,
)
