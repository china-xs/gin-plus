// Package http
// @file      : middleware.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 17:35
// @Description:
package http

import "github.com/gin-gonic/gin"

// Handler defines the handler invoked by Middleware.
type Handler func(c *gin.Context, req interface{}) (any, error)

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(Handler) Handler

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
