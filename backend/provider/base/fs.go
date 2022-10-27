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
	var fi *model.FileInfo
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fi) {
		return *fi, err
	}
	dir, base, parents, rel, isRoot := util.SplitPath(name, path.Join("/", file.ProviderFolder.Name))
	fi = &model.FileInfo{Name: name, Id: file.ProviderFolder.RemoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
	fi.Rel = rel
	fi.Dir = dir
	fi.Base = base
	fi.Parents = parents
	fi.IsRoot = isRoot
	if !isRoot {
		err = findFile(fi)
	}
	if err == nil {
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), &fi)
	}
	return *fi, err
}
