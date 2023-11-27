// Package utils
// @file      : db.log.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/9 17:44
// @Description:
package utils

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
	"time"
)

const _sqlKey = "sql"
const _sqlLog = "sql-log"

var (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

type GormOpts struct {
	Host                   string          `yaml:"host"`     // ip
	Port                   int32           `yaml:"port"`     // 端口
	Database               string          `yaml:"database"` // 表名称
	User                   string          `yaml:"user"`     // 用户名
	Pwd                    string          `yaml:"pwd"`      // 密码
	PreFix                 string          `yaml:"prefix"`   // 表前缀
	MaxIdleConn            int             `yaml:"maxIdleConn"`
	MaxOpenConn            int             `yaml:"maxOpenConn"`     // 最大链接池数
	ConnMaxLifetime        time.Duration   `yaml:"connMaxLifetime"` // 单个链接有效时长
	Level                  logger.LogLevel `yaml:"level"`
	SlowThreshold          time.Duration   `yaml:"slowTime"`               // 慢查询阀值
	SkipDefaultTransaction bool            `yaml:"skipDefaultTransaction"` // true 开启禁用事物，大约 30%+ 性能提升
	SingularTable          bool            `yaml:"singularTable"`
}

// NewGorm gorm v 基础配置 lg gorm 日志文件
func NewGorm(v *viper.Viper, lg logger.Interface) (db *gorm.DB, fc func(), err error) {
	var (
		c      = new(GormOpts)
		config *gorm.Config
	)
	if err = v.UnmarshalKey("db", c); err != nil {
		return
	}

	config = &gorm.Config{
		SkipDefaultTransaction: c.SkipDefaultTransaction,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.PreFix,        // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: c.SingularTable, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		// 数据库配置 统一公用重写日志库即可
		Logger: lg,
	}
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.User,
		c.Pwd,
		c.Host,
		c.Port,
		c.Database,
	)
	if db, err = gorm.Open(mysql.Open(dsn), config); err != nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("db get mysql.DB err%v", err))
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime * time.Second)
	cleanup := func() {
		sqlDB.Close()
	}
	return db, cleanup, nil
}

// NewGormLog 待完成,日志没有添加trace_id
func NewGormLog(v *viper.Viper, l *zap.Logger) (res logger.Interface, err error) {
	var (
		c = new(GormOpts)
	)
	if err = v.UnmarshalKey("db", c); err != nil {
		return
	}
	res = &GormLog{
		log:           l,
		logLevel:      c.Level,
		slowThreshold: c.SlowThreshold,
	}
	return
}

type GormLog struct {
	log           *zap.Logger
	logLevel      logger.LogLevel
	slowThreshold time.Duration // 慢查询阀值
}

func (l *GormLog) WithCtx(ctx context.Context) *zap.Logger {
	return WithCtx(ctx, l.log)
}

func (l *GormLog) LogMode(level logger.LogLevel) logger.Interface {
	nlg := *l
	nlg.logLevel = level
	return &nlg
}
func (l *GormLog) Info(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= logger.Info {
		l.WithCtx(ctx).Info(msg, zap.Any(`data`, data))
	}
}
func (l *GormLog) Warn(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= logger.Warn {
		l.WithCtx(ctx).Warn(msg, zap.Any(`data`, data))
	}
}
func (l *GormLog) Error(ctx context.Context, msg string, data ...any) {
	if l.logLevel >= logger.Error {
		l.WithCtx(ctx).Error(msg, zap.Any(`data`, data))
	}

}

func (l *GormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}
	var fields zap.Field
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			fields = zap.String(_sqlKey, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			fields = zap.String(_sqlKey, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
		l.WithCtx(ctx).Error(_sqlLog, fields)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		if rows == -1 {
			fields = zap.String(_sqlKey, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			fields = zap.String(_sqlKey, fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
		l.WithCtx(ctx).Warn(_sqlLog, fields)
	case l.logLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			fields = zap.String(_sqlKey, fmt.Sprintf(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			fields = zap.String(_sqlKey, fmt.Sprintf(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
		l.WithCtx(ctx).Info(_sqlLog, fields)
	}
	return
}
