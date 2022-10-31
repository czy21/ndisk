package model

type ProviderAccountPO struct {
	BaseEntity[int64]
	TrackEntity
	Endpoint string `gorm:"column:endpoint" json:"endpoint"`
	UserName string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Token    string `gorm:"column:token" json:"token"`
	Kind     string `gorm:"column:kind" json:"kind"`
	PutBuf   int    `gorm:"column:put_buf" json:"putBuf"`
	GetBuf   int    `gorm:"column:get_buf" json:"getBuf"`
}

func (ProviderAccountPO) TableName() string {
	return "provider_account"
}

type ProviderFolderPO struct {
	BaseEntity[int64]
	TrackEntity
	Name              string `gorm:"column:name" json:"name"`
	RemoteName        string `gorm:"column:remote_name" json:"remoteName"`
	ProviderAccountId int64  `gorm:"column:provider_account_id" json:"providerAccountId"`
}

func (ProviderFolderPO) TableName() string {
	return "provider_folder"
}

type ProviderFolderBO struct {
	ProviderFolderPO
	Account ProviderAccountPO `gorm:"foreignKey:Id;references:ProviderAccountId;" json:"account"`
}

type ProviderFolderMeta struct {
	ProviderFolderBO
}
type ProviderFileMeta struct {
	Path    string
	Name    string
	Rel     string
	Base    string
	Dir     string
	Parents []string
	IsRoot  bool
}
type ProviderFile struct {
	Target ProviderFileMeta
	Source ProviderFileMeta

	ProviderFolder ProviderFolderMeta
}
