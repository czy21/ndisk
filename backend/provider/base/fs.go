package base

import (
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
)

func GetFileInfo(ctx context.Context, name string, file model.ProviderFile, findFile func(fileInfo *model.FileInfo) error) (model.FileInfo, error) {
	var err error
	fileInfo := file.FileInfo
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return *fileInfo, err
	}
	if !file.IsRoot {
		err = findFile(fileInfo)
	}
	if err == nil {
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo)
	}
	return *fileInfo, err
}
