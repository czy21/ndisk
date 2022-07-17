package _189

import (
	"github.com/czy21/ndisk/model"
	"strconv"
	"time"
)

type Response struct {
	ResCode   int    `json:"res_code"`
	ResMsg    string `json:"res_message"`
	Success   string `json:"success"`
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
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

type FolderMetaAddRes struct {
	Response
	FolderMeta
}

func (f FileMeta) MapToFileInfo() model.FileInfo {
	return model.FileInfo{
		IsDir:      false,
		ModTime:    time.Time(f.UpdateDate),
		Size:       f.Size,
		RemoteName: strconv.FormatInt(f.Id, 10),
	}
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

type TaskRes struct {
	Response
	TaskId string `json:"taskId"`
}

type FileDownloadUrlRes struct {
	Response
	Url string `json:"fileDownloadUrl"`
}
