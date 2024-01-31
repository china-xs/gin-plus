// Package http
// @file      : codec.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 17:34
// @Description:
package http

import (
	"github.com/china-xs/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/trace"
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

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(c *gin.Context, obj any, err error) {
	var traceId string
	if span := trace.SpanContextFromContext(c.Request.Context()); span.HasTraceID() {
		traceId = span.TraceID().String()
	}
	if err != nil {
		coder := errors.ParseCoder(err)
		var message, reason string
		var code, httpInt int
		if coder.Code() == 1 {
			message = err.Error()
			code = 1
			httpInt = 200
		} else {
			message = coder.String()
			reason = coder.Reference()
			code = coder.Code()
			httpInt = coder.HTTPStatus()
		}
		errCause := errors.Cause(err)
		if e, ok := errCause.(*errors.CustomErr); ok {
			message = e.Error()
			code = e.Code()
		}
		c.JSON(httpInt, map[string]any{
			`msg`:      message,
			`code`:     code,
			`reason`:   reason,
			`trace_id`: traceId,
		})
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		`msg`:      `ok`,
		`code`:     0,
		`data`:     obj,
		`trace_id`: traceId,
	})
}
