package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mediahub/pkg/config"
	"mediahub/services"
	"mediahub/services/shorturl"
	"mediahub/services/shorturl/proto"

	"mediahub/pkg/log"
	"mediahub/pkg/storage"
	"mediahub/pkg/utils"
	"mediahub/pkg/xerror"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	os     storage.Storage
	log    log.ILogger
	config *config.Config
}

func NewController(os storage.Storage, log log.ILogger, cnf *config.Config) *Controller {
	return &Controller{
		os:     os,
		log:    log,
		config: cnf,
	}
}

// Upload 上传文件
func (c *Controller) Upload(ctx *gin.Context) {
	// 1. 获取用户身份与接收文件
	userId := ctx.GetInt64("User.ID")
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	// 2. 读取文件内容到内存
	file, err := fileHeader.Open()
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	// 3. 文件格式校验
	// 校验方式不依赖文件后缀名，而是通过读取文件内容的魔数来判断是否为真实图片
	// 防止用户将.exe或.sh等非图片文件重命名为.jpg上传。
	if !utils.IsImage(bytes.NewReader(content)) {
		err = xerror.NewByMsg("仅支持jpg、png、gif格式")
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	// 4. 生成唯一文件名与路径规划
	md5 := utils.MD5(content)
	filename := fmt.Sprintf("%x%s", md5, path.Ext(fileHeader.Filename))
	filePath := "/public/" + filename
	if userId != 0 {
		filePath = fmt.Sprintf("/%d/%s", userId, filename)
	}
	// 5. 执行上传
	// 将内存中的数据再次包装成流，传递给OSS进行上传，返回的url是文件在OSS中完整的访问地址，即长链接
	url, err := c.os.Upload(bytes.NewReader(content), md5, filePath)
	// 6. 生成短链接
	// 6.1 获取短链接服务连接池
	shortPool := shorturl.NewShortUrlClientPool()
	conn := shortPool.Get()
	defer shortPool.Put(conn)
	// 6.2 初始化gRPC客户端
	client := proto.NewShortUrlClient(conn)
	// 6.3 构造请求
	in := &proto.Url{
		Url:      url,
		UserID:   userId,
		IsPublic: userId == 0,
	}
	// 6.4 设置gRPC上下文
	outGoingCtx := context.Background()
	outGoingCtx = services.AppendBearerTokenToContext(outGoingCtx, c.config.DependOn.ShortUrl.AccessToken)
	// 6.5 调用远程服务
	outUrl, err := client.GetShortUrl(outGoingCtx, in)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"url": outUrl.Url,
	})
	return
}
