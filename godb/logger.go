package godb

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"

	"github.com/happyxhw/pkg/constant"
	"github.com/happyxhw/pkg/cx"
)

type gormLogger struct {
	logger          *zap.Logger
	level           gLogger.LogLevel
	slowThreshold   time.Duration
	sqlLenThreshold int
}

func newLogger(logger *zap.Logger, level string, slowThreshold time.Duration, sqlLenThreshold int) gormLogger {
	ll := gLogger.Warn
	switch strings.ToLower(level) {
	case "info":
		ll = gLogger.Info
	case "warn":
		ll = gLogger.Warn
	case "error":
		ll = gLogger.Error
	case "silent":
		ll = gLogger.Silent
	}
	if sqlLenThreshold == 0 {
		sqlLenThreshold = 1000
	}
	gl := gormLogger{
		logger:          logger,
		level:           ll,
		slowThreshold:   slowThreshold,
		sqlLenThreshold: sqlLenThreshold,
	}
	return gl
}

func (gl gormLogger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	return gormLogger{
		logger:        gl.logger,
		level:         level,
		slowThreshold: gl.slowThreshold,
	}
}

func (gl gormLogger) Info(_ context.Context, s string, i ...interface{}) {
	if gl.level < gLogger.Info {
		return
	}
	gl.logger.Sugar().Infof(s, i...)
}

func (gl gormLogger) Warn(_ context.Context, s string, i ...interface{}) {
	if gl.level < gLogger.Warn {
		return
	}
	gl.logger.Sugar().Warnf(s, i...)
}

func (gl gormLogger) Error(_ context.Context, s string, i ...interface{}) {
	if gl.level < gLogger.Error {
		return
	}
	gl.logger.Sugar().Errorf(s, i...)
}

func (gl gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.level <= 0 {
		return
	}
	elapsed := time.Since(begin)
	reqID := cx.RequestID(ctx)
	switch {
	case gl.level >= gLogger.Error && err != nil && err != gorm.ErrRecordNotFound:
		sql, rows := fc()
		sql = gl.trimSQL(sql)
		gl.logger.Error("[GORM]", zap.Error(err), zap.Int64("elapsed", elapsed.Milliseconds()),
			zap.Int64("rows", rows), zap.String("sql", sql), zap.String(constant.HeaderXRequestID, reqID))
	case gl.level >= gLogger.Warn && gl.slowThreshold != 0 && elapsed > gl.slowThreshold:
		sql, rows := fc()
		sql = gl.trimSQL(sql)
		gl.logger.Warn("[GORM]", zap.Int64("elapsed", elapsed.Milliseconds()),
			zap.Int64("rows", rows), zap.String("sql", sql), zap.String(constant.HeaderXRequestID, reqID))
	case gl.level >= gLogger.Info:
		// 排查 prometheus 的查询
		if cx.FromMetricCtx(ctx) {
			return
		}
		sql, rows := fc()
		sql = gl.trimSQL(sql)
		gl.logger.Info("[GORM]", zap.Int64("elapsed", elapsed.Milliseconds()),
			zap.Int64("rows", rows), zap.String("sql", sql), zap.String(constant.HeaderXRequestID, reqID))
	}
}

func (gl gormLogger) trimSQL(sql string) string {
	if gl.sqlLenThreshold == 0 || len(sql) <= gl.sqlLenThreshold {
		return sql
	}
	return sql[:gl.sqlLenThreshold]
}
