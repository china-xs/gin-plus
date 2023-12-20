// Package http
// @file      : codec.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 17:34
// @Description:
package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*gin.Context, any) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(*gin.Context, any, error)

// DefaultRequestDecoder decodes the request body to object.
func DefaultRequestDecoder(c *gin.Context, v any) (err error) {
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch: // default Content-Type:json
		if err = c.ShouldBindBodyWith(v, binding.JSON); err != nil {
			//err = errors.WrapC(err, http.StatusBadRequest, ``)
			return err
		}
	case http.MethodGet, http.MethodDelete:
		if err = c.ShouldBindQuery(v); err != nil {
			//err = errors.WrapC(err, http.StatusBadRequest, ``)
			return
		}
	}
	return
}

var response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(c *gin.Context, obj any, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			`message`: err.Error(),
			`code`:    http.StatusBadRequest,
		})
		//coder := errors.ParseCoder(err)
		//c.JSON(coder.Code(), map[string]any{
		//	`message`: coder.String(),
		//	`code`:    coder.Code(),
		//})
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		`message`: `ok`,
		`code`:    http.StatusOK,
		`data`:    obj,
	})
}
