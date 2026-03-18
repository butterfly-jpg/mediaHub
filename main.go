package main

import (
	"flag"
	"fmt"
	"mediahub/controller"
	"mediahub/middleware"
	"mediahub/pkg/config"
	"mediahub/pkg/log"
	"mediahub/pkg/storage/cos"
	"mediahub/routers"
	"net/http"

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

	log.SetLevel(cnf.Log.Level)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetPrintCaller(true)

	logger := log.NewLogger()
	logger.SetLevel(cnf.Log.Level)
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetPrintCaller(true)

	// 初始化cos
	os := cos.NewCosStorage(cnf.Cos.BucketUrl, cnf.Cos.SecretId, cnf.Cos.SecretKey, cnf.Cos.CDNDomain)
	c := controller.NewController(os, logger, cnf)

	// 启动应用程序
	gin.SetMode(cnf.Server.Mode)
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.Auth())
	// 注册了一个/health健康检测接口，用于监控存活
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api := r.Group("/api")
	routers.InitRouters(api, c)

	fs := http.FileServer(http.Dir("www"))
	r.NoRoute(func(c *gin.Context) {
		fs.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "www/index.html")
	})

	if err := r.Run(fmt.Sprintf("%s:%d", cnf.Server.IP, cnf.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
