package cos

import (
	"context"
	"encoding/base64"
	"io"
	"mediahub/pkg/storage"
	"mime"
	"net/http"
	"path"

	url2 "net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type cosStorage struct {
	bucketUrl string
	secretId  string
	secretKey string
	cdnDomain string
}

func NewCosStorage(bucketUrl, secretId, secretKey, cdnDomain string) storage.Storage {
	return &cosStorage{
		bucketUrl: bucketUrl,
		secretId:  secretId,
		secretKey: secretKey,
		cdnDomain: cdnDomain,
	}
}

func (s *cosStorage) Upload(r io.Reader, md5Digest []byte, dstPath string) (url string, err error) {
	u, _ := url2.Parse(s.bucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  s.secretId,
			SecretKey: s.secretKey,
		},
	})

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: mime.TypeByExtension(path.Ext(dstPath)),
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{},
	}
	if len(md5Digest) != 0 {
		opt.ObjectPutHeaderOptions.ContentMD5 = base64.StdEncoding.EncodeToString(md5Digest)
	}

	_, err = client.Object.Put(context.Background(), dstPath, r, opt)
	if err != nil {
		return "", err
	}
	url = s.bucketUrl + dstPath
	if s.cdnDomain != "" {
		url = s.cdnDomain + dstPath
	}
	return

}
