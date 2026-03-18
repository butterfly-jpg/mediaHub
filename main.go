package main

import (
	"flag"
	"mediahub/pkg/config"
	"mediahub/pkg/storage/oss"
	"net/http"

	"mediahub/todo"

	"github.com/gin-gonic/gin"
)

var (
	configFile = flag.String("config", "dev.config.yaml", "")
)

func main() {
	// 可获取命令行输入的配置参数
	flag.Parse()
	// 初始化配置文件
	config.InitConfig(*configFile)
	cnf := config.GetConfig()
	// 初始化oss
	os := oss.NewOssStorage(cnf.OSS.BucketName, cnf.OSS.OssAccessKeyID, cnf.OSS.OssAccessKeySecret, cnf.OSS.RegionId)
	
	r := gin.Default()

	// Serve frontend
	r.StaticFS("/static", http.Dir("./static"))
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Todo API
	api := r.Group("/api/v1/todos")
	todo.RegisterRoutes(api)

	r.Run(":8080")
}
