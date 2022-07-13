package model

type ProviderAccountPO struct {
	BaseEntity[int64]
	TrackEntity
	UserName string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Token    string `gorm:"column:token" json:"token"`
	Kind     string `gorm:"column:kind" json:"kind"`
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

type Provider struct {
	ProviderFolderPO
	Account ProviderAccountPO `gorm:"foreignKey:Id;references:ProviderAccountId;" json:"account"`
}
