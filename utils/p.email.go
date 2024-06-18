// Package utils
// @file      : p.email.go
// @author    : xs
// @time      : 2024/6/17 13:59
// @Description: 邮件警告通知
package utils

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"log"
)

type MailEntity struct {
	MainHost  string              `yaml:"mainHost"`
	MailFrom  string              `yaml:"mailFrom"`  // 哨兵、发送邮件账号
	MailPwd   string              `yaml:"mailPwd"`   // 发送邮件密码
	MailPort  int                 `yaml:"mailPort"`  // 发送邮件端口
	MailGroup map[string][]string `yaml:"mailGroup"` // 邮件组
	*Log
}

func NewMail(v *viper.Viper, logger *zap.Logger) *MailEntity {
	var err error
	var entity = new(MailEntity)
	if err = v.UnmarshalKey("sentry", entity); err != nil {
		panic(fmt.Sprintf(`init email-sentry-err: %s\n`, err.Error()))
	}
	entity.Log = NewWLog(logger)
	return entity
}

// GetMailGroup 获取邮件组下所有邮箱 all 返回所有, map[] 不可以设置 all
func (this *MailEntity) GetMailGroup(key string) (result []string) {
	for k, strings := range this.MailGroup {
		if key == `all` || k == key {
			result = append(result, strings...)
		}
	}
	return
}

func (this *MailEntity) Send(ctx context.Context, to []string, title string, content string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", this.MailFrom)
	m.SetHeaders(map[string][]string{"To": to})
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)
	d := gomail.NewDialer(this.MainHost, this.MailPort, this.MailFrom, this.MailPwd) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err = d.DialAndSend(m); err != nil {
		this.WithCtx(ctx).Warn(`email-api`, zap.Error(err))
		log.Printf("发送失败:%s", err)
		return
	}
	return
}
