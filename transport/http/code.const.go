// Package http
// @file      : code.con.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/12/21 15:20
// @Description:
package http

const (
	ErrSucc         = 100000         // 成功
	ErrBodyEOF  int = iota + ErrSucc // post 请求参数为空
	ErrValidate                      //  请求参数校验异常

)
