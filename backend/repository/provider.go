package repository

import (
	"github.com/czy21/cloud-disk-sync/model"
)

type Provider struct {
}

//func (Environment) SelectListBy(query model.EnvironmentQuery) []model.EnvironmentPO {
//	d := dbClient.Model(&model.EnvironmentPO{})
//	if len(query.Name) > 0 {
//		d.Where(" name like ? ", query.Name+"%")
//	}
//	var list []model.EnvironmentPO
//	d.Find(&list)
//	return list
//}

func (Provider) InsertOneForAccount(po model.ProviderAccountPO) {
	dbClient.Create(&po)
}

func (Provider) InsertOneForFolder(po model.ProviderFolderPO) {
	dbClient.Create(&po)
}

func (Provider) SelectList() []model.ProviderMeta {
	d := dbClient.Model(&model.ProviderMeta{})
	d.Preload("Account")
	var list []model.ProviderMeta
	d.Find(&list)
	return list
}
