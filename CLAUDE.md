# Go 项目开发规范

## 项目概览

- **模块名**: `mediahub`
- **Go 版本**: 1.25+
- **Web 框架**: Gin
- **配置管理**: Viper（支持热更新）
- **数据库**: MySQL（标准库 `database/sql`）、Redis（`go-redis/v9`）
- **文件存储**: 阿里云 OSS
- **日志**: Logrus + Lumberjack（按日期轮转）
- **错误处理**: 自定义 `pkg/xerror` 包

## 目录结构

```
mediahub/
├── main.go                 # 入口：初始化配置、依赖，注册路由
├── middleware/             # HTTP 中间件（Auth、CORS 等）
├── controller/             # HTTP 处理器（较大模块独立放这里）
├── services/               # 业务逻辑层、gRPC 客户端封装
├── pkg/
│   ├── config/             # 配置结构体与初始化
│   ├── db/
│   │   ├── mysql/          # MySQL 初始化与获取
│   │   └── redis/          # Redis 初始化、Key 前缀
│   ├── log/                # 日志初始化与 Hook
│   ├── storage/            # 存储抽象接口
│   │   └── oss/            # 阿里云 OSS 实现
│   ├── constants/          # 全局常量
│   └── xerror/             # 自定义错误类型
└── {feature}/              # 小模块可按业务聚合（如 todo/）
    ├── handler.go          # HTTP 路由注册与处理函数
    ├── service.go          # 业务逻辑
    └── model.go            # 数据模型
```

## 命名规范

- **包名**: 全小写，单词，与目录名一致（如 `package mysql`）
- **文件名**: 全小写，下划线分隔（如 `rotate_writer.go`）
- **变量/函数**: 导出用 PascalCase，未导出用 camelCase
- **常量**: PascalCase 或全大写 + 下划线（按语义选择）
- **接口**: 以功能动词或 `-er` 结尾（如 `Storage`、`Writer`）
- **结构体字段 JSON tag**: snake_case（如 `json:"avatar_url"`）
- **配置 mapstructure tag**: camelCase（如 `mapstructure:"bucketName"`）

## 代码规范

### 包初始化模式

全局单例用包级私有变量 + Init/Get 函数：

```go
var db *sql.DB

func InitMysql(cnf *config.Config) { ... }
func GetDB() *sql.DB { return db }
```

### 错误处理

- 优先使用 `pkg/xerror` 包创建语义化错误
- 初始化阶段的致命错误使用 `panic` 或 `log.Fatal`
- HTTP 处理层统一用 `c.JSON` 返回错误，不要 panic

```go
// 带错误码
xerror.NewByCode("ERR_001", "资源不存在")
// 包装原始错误
xerror.NewByErr(err)
```

### HTTP 响应格式

统一使用 `gin.H` 返回 JSON：

```go
// 成功
c.JSON(http.StatusOK, gin.H{"data": result})
// 列表
c.JSON(http.StatusOK, gin.H{"data": list, "total": len(list)})
// 错误
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 创建
c.JSON(http.StatusCreated, gin.H{"data": created})
```

### 请求绑定

使用 `ShouldBindJSON` 并在失败时立即返回：

```go
var req createRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

### 路由注册

每个模块提供 `RegisterRoutes(rg *gin.RouterGroup)` 函数，在 `main.go` 中挂载：

```go
api := r.Group("/api/v1/short-links")
shortlink.RegisterRoutes(api)
```

### 中间件

返回 `gin.HandlerFunc` 的闭包：

```go
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ...
        c.Next()
    }
}
```

### 配置访问

通过 `config.GetConfig()` 获取全局配置，不要在函数签名中传递 `*Config`，直接在函数体内调用：

```go
conf := config.GetConfig()
```

### 接口与实现分离

存储等基础设施定义接口，具体实现放子包：

```go
// pkg/storage/storage.go
type Storage interface {
    Upload(ctx context.Context, key string, r io.Reader) error
}

// pkg/storage/oss/oss.go
type OssStorage struct { ... }
func (o *OssStorage) Upload(...) error { ... }
```

## 依赖使用规范

| 场景 | 使用方式 |
|------|----------|
| HTTP 路由 | `github.com/gin-gonic/gin` |
| 配置读取 | `github.com/spf13/viper` + `mapstructure` tag |
| 日志 | `github.com/sirupsen/logrus` |
| Redis | `github.com/redis/go-redis/v9` |
| MySQL | 标准库 `database/sql` + `_ "github.com/go-sql-driver/mysql"` |
| OSS | `github.com/aliyun/alibabacloud-oss-go-sdk-v2` |
| 请求校验 | `binding:"required"` tag（gin 内置 validator） |

## 新功能开发流程

1. 在 `pkg/config/config.go` 的 `Config` 结构体中添加配置字段
2. 在对应目录创建 `model.go`（数据结构）、`service.go`（业务逻辑）、`handler.go`（HTTP 层）
3. 在 `handler.go` 中实现 `RegisterRoutes`
4. 在 `main.go` 中初始化依赖并注册路由
5. 如需新的 pkg 基础设施，在 `pkg/` 下创建对应子包

## 注意事项

- 不要在 handler 层写业务逻辑，保持 handler 薄，业务放 service
- Redis Key 前缀统一在 `pkg/db/redis/prefix.go` 中定义常量
- 全局常量放 `pkg/constants/constants.go`
- 不要直接使用 `fmt.Println` 打印日志，统一用 logrus
- Context 传递用 `context.Context`，gRPC 调用用 `services.AppendBearerTokenToContext`