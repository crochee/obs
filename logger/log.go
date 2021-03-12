// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package logger
package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    DefaultLogSizeM, //单个日志文件最大MaxSize*M大小 // megabytes
		MaxAge:     MaxLogDays,      //days
		MaxBackups: DefaultMaxZip,   //备份数量
		Compress:   false,           //不压缩
		LocalTime:  true,            //备份名采用本地时间
	}
}

// NewLogger 初始化日志对象
//
// @param: path 日志路径
// @param: level 日志等级
func NewLogger(path, level string) *Logger {
	logger := &Logger{
		Logger: NewZap(level, zapcore.NewJSONEncoder, SetLoggerWriter(path)),
	}
	logger.LoggerSugar = logger.Logger.Sugar()
	return logger
}

type Logger struct {
	Logger      *zap.Logger
	LoggerSugar *zap.SugaredLogger
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l == nil {
		return
	}
	l.LoggerSugar.Debugf(format, v...)
}

// Debug 打印Debug信息
//
// @param: message 格式信息
func (l *Logger) Debug(message string) {
	if l == nil {
		return
	}
	l.Logger.Debug(message)
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Infof(format string, v ...interface{}) {
	if l == nil {
		return
	}
	l.LoggerSugar.Infof(format, v...)
}

// Info 打印Info信息
//
// @param: message 格式信息
func (l *Logger) Info(message string) {
	if l == nil {
		return
	}
	l.Logger.Info(message)
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l == nil {
		return
	}
	l.LoggerSugar.Errorf(format, v...)
}

// Error 打印Error信息
//
// @param: message 信息
func (l *Logger) Error(message string) {
	if l == nil {
		return
	}
	l.Logger.Error(message)
}
