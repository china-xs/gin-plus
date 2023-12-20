// Package http
// @file      : mw.logger.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/12/20 14:48
// @Description: 记录系统日志
package http

import (
	"fmt"
	"github.com/china-xs/gin-plus/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func MwLogger(lg *zap.Logger) Middleware {
	l := utils.NewWLog(lg)
	return func(handler Handler) Handler {
		return func(c *gin.Context, req interface{}) (reply any, err error) {
			startTime := time.Now()
			reply, err = handler(c, req)
			var body string
			// 记录body 数据
			if bodyBytes, ok := c.Get(gin.BodyBytesKey); ok {
				body = string(bodyBytes.([]byte))
			}
			var fields = make([]zap.Field, 9)
			fields[0] = zap.String("url", c.Request.URL.String())
			fields[1] = zap.String("method", c.Request.Method)
			fields[2] = zap.String("body", body)
			fields[3] = zap.String("host", c.Request.Host)
			fields[4] = zap.String("ipv4", c.ClientIP())
			fields[5] = zap.String("latency", time.Since(startTime).String())
			fields[6] = zap.String("args", extractArgs(req))
			fields[7] = zap.String("reply", extractArgs(reply))
			fields[8] = zap.Error(err)
			//fields[9] = zap.String("reason", reason)
			// 当前仅记录 到api曾的出入参数，如需独立记录额外参数，请独立配置
			l.WithCtx(c.Request.Context()).Info(`req-log`, fields...)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
