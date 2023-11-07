// Package ginplus
// @file      : options.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 16:54
// @Description:
package ginplus

import (
	"context"
	"os"
	"time"
)

type Option func(o *options)

type options struct {
	ctx     context.Context
	sigs    []os.Signal
	timeout time.Duration
	servers []Server
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Timeout with server timeout time
func Timeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}

// Servers with start & stop server.
func Servers(srv ...Server) Option {
	return func(o *options) { o.servers = srv }
}
