// Package utils
// @file      : p.redis.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/10 09:08
// @Description: redis config
package utils

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisOpt struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"` // user 没有account 可不填写
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"poolSize"`
}

func NewRedis(v *viper.Viper) (rdb *redis.Client, err error) {
	var (
		c   = new(RedisOpt)
		opt *redis.Options
	)
	if err = v.UnmarshalKey("redis", c); err != nil {
		return
	}
	opt = &redis.Options{
		Addr:     c.Addr, //localhost:6379
		Username: c.User,
		Password: c.Password,
		DB:       c.DB,
		PoolSize: c.PoolSize,
	}
	rdb = redis.NewClient(opt)
	return
}
