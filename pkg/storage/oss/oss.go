package oss

import (
	"context"
	"mediahub/pkg/storage"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type ossStorage struct {
	ossAccessKeyId     string
	ossAccessKeySecret string
	regionId           string
	bucketName         string
	objName            string
}

func NewOssStorage(bucketName, ossAccessKeyId, ossAccessKeySecret, regionId string) storage.Storage {
	return &ossStorage{
		ossAccessKeyId:     ossAccessKeyId,
		ossAccessKeySecret: ossAccessKeySecret,
		regionId:           regionId,
		bucketName:         bucketName,
	}
}

func (o *ossStorage) Upload(content string) error {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(o.regionId)
	// 创建oss客户端
	client := oss.NewClient(cfg)
	// 定义要上传的内容
	body := strings.NewReader(content)
	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(o.bucketName),
		Key:    oss.Ptr(o.objName),
		Body:   body,
	}
	// 发送上传对象的请求
	_, err := client.PutObject(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}
