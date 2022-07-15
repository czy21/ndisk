package _189

import (
	"github.com/czy21/cloud-disk-sync/model"
)

type Response struct {
	Code    int    `json:"res_code"`
	Message string `json:"res_message"`
}

type BaseTrackModel[TID any] struct {
	Id         TID                `json:"id"`
	CreateDate model.StandardTime `json:"createDate"`
	UpdateDate model.StandardTime `json:"lastOpTime"`
}

type FolderMeta struct {
	BaseTrackModel[int64]
	Name string `json:"name"`
}

type FileMeta struct {
	BaseTrackModel[int64]
	Name string `json:"name"`
	Size int64  `json:"size"`
}
type FileListAO struct {
	Count   int          `json:"count"`
	Files   []FileMeta   `json:"fileList"`
	Folders []FolderMeta `json:"folderList"`
}
type FileListAORes struct {
	Response
	FileListAO FileListAO `json:"fileListAO"`
}
