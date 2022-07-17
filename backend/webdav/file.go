package webdav

import (
	"github.com/czy21/ndisk/model"
	"github.com/spf13/viper"
	"io/fs"
	"os"
)

type File struct {
	Name string
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
	var fileInfos []fs.FileInfo
	for _, t := range providerMetas {
		fileInfos = append(fileInfos,
			model.FileInfoProxy{
				FileInfo: model.FileInfo{
					IsDir: true,
					Name:  t.Name,
				},
			})
	}
	ds, _ := os.ReadDir(viper.GetString("data.dav"))
	for _, t := range ds {
		if t.IsDir() {
			for _, pm := range providerMetas {
				if t.Name() != pm.Name {
					fileInfos = append(fileInfos,
						model.FileInfoProxy{
							FileInfo: model.FileInfo{
								Name:  t.Name(),
								IsDir: true,
							},
						})
				}
			}
		}
	}
	return fileInfos, nil
}

func (f File) Stat() (fs.FileInfo, error) {
	return model.FileInfoProxy{FileInfo: model.FileInfo{IsDir: true}}, nil
}

func (f File) Write(p []byte) (n int, err error) {
	panic("implement me")
}
