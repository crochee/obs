// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/6

package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"obs/controller/file"

	"obs/controller/bucket"
	_ "obs/docs"
	"obs/middleware"
)

// @title obs Swagger API
// @version 1.0
// @description This is a obs server.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Auth-Token

// GinRun gin router
//
// @Success *gin.Engine gin router
func GinRun() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CrossDomain, middleware.TraceId, middleware.Recovery, middleware.Log)
	if gin.Mode() != gin.ReleaseMode {
		// swagger
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		// 增加性能测试
		pprof.Register(router)
	}

	v1Router := router.Group("/v1")

	{
		// bucket
		v1Router.POST("/bucket", bucket.CreateBucket)
		v1Router.HEAD("/bucket/:bucket_id", bucket.HeadBucket)
		v1Router.DELETE("/bucket/:bucket_id", bucket.DeleteBucket)

		// file
		fileRouter := v1Router.Group("/bucket/:bucket_id")
		{
			fileRouter.POST("/file", file.UploadFile)
			fileRouter.DELETE("/file/:file_id", file.DeleteFile)
			fileRouter.HEAD("/file/:file_id", file.SignFile)
			fileRouter.GET("/file", file.DownloadFile)
			fileRouter.GET("/files", file.FileList)
		}
	}

	fileRouter := v1Router.Group("/file")
	{
		fileRouter.POST("/:bucket_name", file.UploadFile)
		fileRouter.DELETE("/:bucket_name", file.DeleteFile)
		fileRouter.HEAD("/:bucket_name", file.SignFile)
		fileRouter.GET("/:bucket_name", file.DownloadFile)
	}
	return router
}
