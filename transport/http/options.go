// Package http
// @file      : options.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/8 11:32
// @Description:
package http

import "github.com/gin-gonic/gin"

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// IsDebug open gin debug model
func IsDebug() ServerOption {
	return func(o *Server) {
		o.isDebug = true
	}
}

// Address with server address.
func Address(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

// RequestDecoder with request decoder.
func RequestDecoder(dec DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.dec = dec
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(enc EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.enc = enc
	}
}

// Middlewares set middleware
func Middlewares(ms ...Middleware) ServerOption {
	return func(o *Server) {
		o.ms = ms
	}
}

// Filter with HTTP middleware option.
func Filter(filters ...gin.HandlerFunc) ServerOption {
	return func(o *Server) {
		o.filters = filters
	}
}
