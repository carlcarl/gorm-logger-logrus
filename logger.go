package gormloggerlogrus

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Options struct {
	Logger *logrus.Logger

	LogLevel                  gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	FileWithLineNumField      string
}

type Logger struct {
	Options
}

func New(opts Options) *Logger {
	l := &Logger{Options: opts}
	if l.Logger == nil {
		l.Logger = logrus.New()
	}
	if l.LogLevel == 0 {
		l.LogLevel = gormlogger.Silent
	}

	return l
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Logger.WithContext(ctx).Infof(s, args...)
	}
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Logger.WithContext(ctx).Warnf(s, args...)
	}
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Logger.WithContext(ctx).Errorf(s, args...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	fields := logrus.Fields{}
	if l.FileWithLineNumField != "" {
		fields[l.FileWithLineNumField] = utils.FileWithLineNum()
	}

	sql, rows := fc()
	if rows == -1 {
		fields["rows_affected"] = "-"
	} else {
		fields["rows_affected"] = rows
	}

	elapsed := time.Since(begin)
	fields["elapsed"] = elapsed

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) && l.LogLevel >= gormlogger.Error:
		l.Logger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		l.Logger.WithContext(ctx).WithFields(fields).Warnf("SLOW SQL >= %v, l.SlowThreshold %s [%s]", l.SlowThreshold, sql, elapsed)
	case l.LogLevel == gormlogger.Info:
		l.Logger.WithContext(ctx).WithFields(fields).Infof("%s [%s]", sql, elapsed)
	}
}
