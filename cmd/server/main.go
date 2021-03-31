// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/asim/go-micro/v3"

	"obs/cmd"
	"obs/config"
	"obs/cron"
	"obs/logger"
	"obs/model/db"
	"obs/pprof"
	"obs/router"
)

// todo 集成go-micro
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消
	// 初始化配置
	config.InitConfig(os.Getenv("config"))
	// 初始化系统日志
	logger.InitSystemLogger(config.Cfg.ServiceConfig.ServiceInfo.LogPath,
		config.Cfg.ServiceConfig.ServiceInfo.LogLevel)

	var srv *http.Server
	service := micro.NewService(
		micro.Context(ctx),
		micro.Name(cmd.ServiceName),
		micro.Version(cmd.Version),
		micro.Address(":8160"),
		micro.HandleSignal(true),
		micro.Profile(pprof.NewPprof("", "")),
		micro.BeforeStart(func() error {
			// cron 初始化
			cron.Setup()
			// 数据库初始化
			if err := db.Setup(ctx); err != nil {
				logger.Fatalf(err.Error())
			}
			// 初始化请求日志
			requestLog := logger.NewLogger(config.Cfg.ServiceConfig.ServiceInfo.LogPath,
				config.Cfg.ServiceConfig.ServiceInfo.LogLevel)
			// http服务对象
			srv = &http.Server{
				Addr:    ":8150",
				Handler: router.GinRun(),
				ConnContext: func(ctx context.Context, c net.Conn) context.Context {
					return logger.With(ctx, requestLog)
				},
				BaseContext: func(listener net.Listener) context.Context {
					return ctx
				},
			}
			return nil
		}),
		micro.AfterStart(func() error {
			go func() {
				logger.Info("obs running...")
				if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
					logger.Errorf(err.Error())
				}
			}()
			return nil
		}),
		micro.BeforeStop(func() error {
			newCtx, newCancel := context.WithTimeout(ctx, 15*time.Second)
			defer newCancel()
			// The context is used to inform the server it has 5 seconds to finish
			// the request it is currently handling
			if err := srv.Shutdown(newCtx); err != nil {
				logger.Errorf("Server forced to shutdown:%v", err)
			}
			return nil
		}),
		micro.AfterStop(func() error {
			// 关闭定时器
			cron.New().Stop()
			// 数据库关闭连接池
			if err := db.Close(); err != nil {
				logger.Errorf("db forced to shutdown:%v", err)
			}
			logger.Exit("obs server exit!")
			return nil
		}),
	)
	service.Server().Handle()
	if err := service.Server().Init(); err != nil {
		logger.Fatal(err.Error())
	}
	if err := service.Run(); err != nil {
		logger.Fatal(err.Error())
	}
}