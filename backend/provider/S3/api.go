package S3

import (
	"github.com/czy21/ndisk/model"
	"github.com/minio/minio-go/v6"
	"path"
	"strings"
	"sync"
)

type API struct {
	File model.ProviderFile
}

var clientMap sync.Map

func (a API) GetClient() (*minio.Client, error) {
	account := a.File.ProviderFolder.Account
	var client *minio.Client
	var err error
	key := account.Endpoint + ":" + account.UserName
	if client, ok := clientMap.Load(key); ok {
		return client.(*minio.Client), err
	}
	client, err = minio.New(account.Endpoint, account.UserName, account.Password, false)
	clientMap.Store(key, client)
	return client, err
}

func (a API) GetObjects(bucketName string, objectPrefix string) (objectInfos []minio.ObjectInfo, err error) {
	objectPrefix = strings.TrimPrefix(path.Join(objectPrefix), "/") + "/"
	client, err := a.GetClient()
	doneCh := make(chan struct{})
	defer close(doneCh)
	for t := range client.ListObjects(bucketName, objectPrefix, false, doneCh) {
		if objectPrefix != t.Key {
			objectInfos = append(objectInfos, t)
		}
	}
	return objectInfos, err
}
