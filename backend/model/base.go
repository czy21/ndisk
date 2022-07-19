package model

import "time"

type PageModel struct {
	PageIndex int   `json:"pageIndex"`
	PageSize  int   `json:"pageSize"`
	Total     int64 `json:"total"`
}

type PageResult[T any] struct {
	List []T       `json:"list"`
	Page PageModel `json:"page,omitempty"`
}

type BaseEntity[TID any] struct {
	Id TID `gorm:"column:id" json:"id"`
}

type TrackEntity struct {
	CreateTime *time.Time `gorm:"column:create_time;default:null" json:"createTime"`
	CreateUser int64      `gorm:"column:create_user" json:"createUser"`
	UpdateTime *time.Time `gorm:"column:update_time;default:null" json:"updateTime"`
	UpdateUser int64      `gorm:"column:update_user" json:"updateUser"`
}

type BaseQuery[TID any] struct {
	Name       string `json:"name"`
	ServiceUrl string `json:"serviceUrl"`
}

type SimpleItemModel[T any] struct {
	Label    string                 `json:"label"`
	Value    T                      `json:"value"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
	Children []SimpleItemModel[T]   `json:"children,omitempty"`
}
