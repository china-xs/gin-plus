// Package http
// @file      : http.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 17:30
// @Description:
package http

import (
	"github.com/gin-gonic/gin"
)

const OperationKey = "operation"

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

type Server struct {
	*gin.Engine
	addr   string             // default 0.0.0.0:8080
	dec    DecodeRequestFunc  // 请求参数绑定结构
	enc    EncodeResponseFunc // 定义返回结构
	ms     []Middleware       // 全局中间价
	filter []gin.HandlerFunc  // gin 全局中间件， 执行比ms 早
}

func (this Server) Start() (err error) {

	return
}

func (this Server) Stop() (err error) {
	return
}

func (this Server) Middleware() {

}
