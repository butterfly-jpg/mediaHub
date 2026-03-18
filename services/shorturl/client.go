package shorturl

import (
	"mediahub/pkg/config"
	grpc_client_pool "mediahub/pkg/grpc-client-pool"
	"mediahub/pkg/log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var pool grpc_client_pool.ClientPool
var once sync.Once

func NewShortUrlClientPool() grpc_client_pool.ClientPool {
	var err error
	if pool != nil {
		return pool
	}
	once.Do(func() {
		cnf := config.GetConfig()
		pool, err = grpc_client_pool.NewPool(cnf.DependOn.ShortUrl.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Error(err)
		}
	})
	return pool
}
