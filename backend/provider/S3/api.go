package S3

import (
	"github.com/czy21/ndisk/model"
	"github.com/minio/minio-go/v6"
	"strings"
)

type API struct {
	File model.ProviderFile
}

func (a API) GetClient() (*minio.Client, error) {
	account := a.File.ProviderFolder.Account
	client, err := minio.New(account.Endpoint, account.UserName, account.Password, false)
	return client, err
}

func (a API) GetObjects(bucketName string, objectPrefix string) (objectInfos []minio.ObjectInfo, err error) {
	objectPrefix = strings.TrimPrefix(objectPrefix, "/")
	client, err := a.GetClient()
	doneCh := make(chan struct{})
	defer close(doneCh)
	for t := range client.ListObjects(bucketName, objectPrefix, false, doneCh) {
		objectInfos = append(objectInfos, t)
	}
	return objectInfos, err
}
