package S3

import (
	"fmt"
	"github.com/minio/minio-go/v6"
	"testing"
)

func TestListObjects(t *testing.T) {
	client, _ := minio.New("", "", "", false)
	doneCh := make(chan struct{})
	defer close(doneCh)
	for t := range client.ListObjects("dsm", "", false, doneCh) {
		fmt.Println(t.Key, t.LastModified)
	}
}
