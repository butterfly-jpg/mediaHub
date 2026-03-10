package services

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// AppendBearerTokenToContext 将BearerToken追加到ctx
func AppendBearerTokenToContext(ctx context.Context, accessToken string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+accessToken)
	return metadata.NewOutgoingContext(ctx, md)
}
