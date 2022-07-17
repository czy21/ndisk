package _189

import (
	"context"
	"github.com/czy21/ndisk/model"
	"io/fs"
)

type File struct {
	Name               string
	ProviderFolderMeta model.ProviderFolderMeta
	Context            context.Context
}

func (f File) Close() error {
	return nil
}

func (f File) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	fileInfo, _ := getFileInfo(f.Name, f.ProviderFolderMeta.RemoteName, f.ProviderFolderMeta)
	folder, err := API{}.getFolderById(fileInfo.RemoteName)
	var fileInfos []fs.FileInfo
	for _, t := range folder.Files {
		fileInfos = append(fileInfos, model.FileInfoProxy{
			FileInfo: model.FileInfo{
				Name: t.Name,
			},
		})
	}
	for _, t := range folder.Folders {
		fileInfos = append(fileInfos, model.FileInfoProxy{
			FileInfo: model.FileInfo{
				Name:  t.Name,
				IsDir: true,
			},
		})
	}
	return fileInfos, err
}

func (f File) Stat() (fs.FileInfo, error) {
	fileInfo, _ := getFileInfo(f.Name, f.ProviderFolderMeta.RemoteName, f.ProviderFolderMeta)
	return model.FileInfoProxy{FileInfo: fileInfo}, nil
}

func (f File) Write(p []byte) (n int, err error) {
	panic("implement me")
}
