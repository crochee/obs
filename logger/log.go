// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package logger
package logger

import (
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"obs/model/mongox"
)

const (
	DefaultLogSizeM int = 20
	DefaultMaxZip   int = 50
	MaxLogDays      int = 30
)

// SetLoggerWriter
func SetLoggerWriter(path string) io.Writer {
	if path == "" {
		return os.Stdout
	}
	if strings.HasPrefix(path, "mongodb://") {
		client, err := mongox.Setup(path)
		if err != nil {
			panic(err)
		}
		return &MongoLogger{client: client}
	}
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    DefaultLogSizeM, //单个日志文件最大MaxSize*M大小 // megabytes
		MaxAge:     MaxLogDays,      //days
		MaxBackups: DefaultMaxZip,   //备份数量
		Compress:   false,           //不压缩
		LocalTime:  true,            //备份名采用本地时间
	}
}

type Option struct {
	Path  string
	Level string
	Skip  int
}

// NewLogger 初始化日志对象
//
// @param: path 日志路径
// @param: level 日志等级
func NewLogger(opts ...func(*Option)) *Logger {
	l := &Logger{
		Option: Option{
			Path:  "",
			Level: "INFO",
			Skip:  1,
		},
	}
	for _, opt := range opts {
		opt(&l.Option)
	}
	var encode encoder
	if l.Option.Path == "" {
		encode = zapcore.NewConsoleEncoder
	} else {
		encode = zapcore.NewJSONEncoder
	}
	l.Logger = newZap(l.Option.Level, encode, l.Option.Skip, SetLoggerWriter(l.Option.Path))
	l.LoggerSugar = l.Logger.Sugar()

	return l
}

type Logger struct {
	Logger      *zap.Logger
	LoggerSugar *zap.SugaredLogger
	Option
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.LoggerSugar.Debugf(format, v...)
}

// Debug 打印Debug信息
//
// @param: message 格式信息
func (l *Logger) Debug(message string) {
	l.Logger.Debug(message)
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Infof(format string, v ...interface{}) {
	l.LoggerSugar.Infof(format, v...)
}

// Info 打印Info信息
//
// @param: message 格式信息
func (l *Logger) Info(message string) {
	l.Logger.Info(message)
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.LoggerSugar.Errorf(format, v...)
}

// Error 打印Error信息
//
// @param: message 信息
func (l *Logger) Error(message string) {
	l.Logger.Error(message)
}

// Fatalf 打印Fatalf信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.LoggerSugar.Errorf(format, v...)
}

// Fatal 打印Fatal信息
//
// @param: message 信息
func (l *Logger) Fatal(message string) {
	l.Logger.Error(message)
}

func (l *Logger) Sync() error {
	err := l.Logger.Sync()
	err1 := l.LoggerSugar.Sync()
	if err1 != nil {
		err = err1
	}
	return err
}
