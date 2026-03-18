package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"encoding/json"
	"mediahub/pkg/config"

	"github.com/gin-gonic/gin"
)

// Auth 授权验证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头中获取token
		token := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")
		// 2. 未携带token不会进行身份验证，允许访客模式
		if token == "" {
			c.Next()
			return
		}
		// 3. 身份验证
		auth, err := checkAuth(token)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		if auth == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 4. 正常通过身份验证，将相关信息通过上下文传递到后续逻辑
		c.Set("User.ID", auth.ID)
		c.Set("User.Name", auth.Name)
		c.Set("User.AvatarUrl", auth.AvatarUrl)
		c.Next()
	}
}

var httpClient = &http.Client{}

// userInfo 用户信息结构体
type userInfo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

// checkAuth 检查用户信息是否正确
func checkAuth(token string) (*userInfo, error) {
	conf := config.GetConfig()
	path := "/api/v1/login/check/auth"
	url := fmt.Sprintf("%s%s?access_token=%s", conf.DependOn.User.Address, path, token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 401 {
		return nil, errors.New("token invalid")
	}
	if res.StatusCode == 500 {
		return nil, errors.New("服务器内部错误")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	user := &userInfo{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
