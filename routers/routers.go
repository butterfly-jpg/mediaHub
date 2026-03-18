package routers

import (
	"mediahub/controller"

	"github.com/gin-gonic/gin"
)

// InitRouters 初始化路由
func InitRouters(api *gin.RouterGroup, c *controller.Controller) {
	v1 := api.Group("/v1")
	v1.POST("/file/upload", c.Upload) // /api/v1/file/upload 会调用上传方法
	v1.GET("/home", c.Home)           // /api/v1/home 会调用返回主页的方法
}
