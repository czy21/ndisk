package base

import (
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/util"
	"golang.org/x/net/context"
	"path"
)

func GetFileInfo(ctx context.Context, name string, file model.ProviderFile, findFile func(fileInfo *model.FileInfo) error) (model.FileInfo, error) {
	var err error
	fileInfo := &model.FileInfo{Name: name, Id: file.ProviderFolder.RemoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return *fileInfo, err
	}
	dir, baseName, dirNames, rel, isRoot := util.SplitPath(name, path.Join("/", file.ProviderFolder.Name))
	fileInfo.Dir = dir
	fileInfo.Base = baseName
	fileInfo.Parents = dirNames
	fileInfo.Rel = rel
	fileInfo.IsRoot = isRoot
	if !isRoot {
		err = findFile(fileInfo)
	}
	if err == nil {
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo)
	}
	return *fileInfo, err
}
