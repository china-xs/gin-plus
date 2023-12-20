// Package http
// @file      : mw.validate.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/12/20 15:14
// @Description:
package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type validator interface {
	Validate() error
}

// MwValidator is a validator middleware.
func MwValidator() Middleware {
	return func(handler Handler) Handler {
		return func(c *gin.Context, req interface{}) (reply any, err error) {
			if v, ok := req.(validator); ok {
				if err = v.Validate(); err != nil {
					return nil, errors.WithMessage(err, `Bad Request`)
				}
			}
			return handler(c, req)
		}
	}
}
