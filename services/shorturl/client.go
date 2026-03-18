package shorturl

import (
	"context"
	"mediahub/pkg/config"
	"mediahub/services"
	"mediahub/services/shorturl/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.ShortUrlClient

// Init 初始化 ShortUrl gRPC 客户端，应在 main.go 中调用。
func Init() {
	conf := config.GetConfig()
	conn, err := grpc.NewClient(
		conf.DependOn.ShortUrl.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic("shorturl: failed to connect gRPC server: " + err.Error())
	}
	client = proto.NewShortUrlClient(conn)
}

// GetShortUrl 根据原始 URL 生成短链接。
func GetShortUrl(ctx context.Context, originalUrl string) (string, error) {
	conf := config.GetConfig()
	ctx = services.AppendBearerTokenToContext(ctx, conf.DependOn.ShortUrl.AccessToken)

	resp, err := client.GetShortUrl(ctx, &proto.GetShortUrlRequest{
		OriginalUrl: originalUrl,
	})
	if err != nil {
		return "", err
	}
	return resp.ShortUrl, nil
}

// GetOriginalUrl 根据短链接还原原始 URL。
func GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	conf := config.GetConfig()
	ctx = services.AppendBearerTokenToContext(ctx, conf.DependOn.ShortUrl.AccessToken)

	resp, err := client.GetOriginalUrl(ctx, &proto.GetOriginalUrlRequest{
		ShortUrl: shortUrl,
	})
	if err != nil {
		return "", err
	}
	return resp.OriginalUrl, nil
}
