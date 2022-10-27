package S3

import (
	"github.com/czy21/ndisk/model"
	"github.com/minio/minio-go/v6"
	"io"
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
	client, err := a.GetClient()
	if err != nil {
		return nil, err
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	for t := range client.ListObjects(bucketName, objectPrefix, false, doneCh) {
		objectInfos = append(objectInfos, t)
	}
	return objectInfos, err
}

func (a API) GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	client, err := a.GetClient()
	if err != nil {
		return nil, err
	}
	return client.GetObject(bucketName, objectName, opts)
}

func (a API) PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (n int64, err error) {
	client, err := a.GetClient()
	if err != nil {
		return 0, err
	}
	return client.PutObject(bucketName, objectName, reader, objectSize, opts)
}

func (a API) StatObject(bucketName string, objectName string) (objectInfo minio.ObjectInfo, err error) {
	client, err := a.GetClient()
	if err != nil {
		return objectInfo, err
	}
	return client.StatObject(bucketName, objectName, minio.StatObjectOptions{})
}
