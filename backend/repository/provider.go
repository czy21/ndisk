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

func (Provider) SelectList() []model.ProviderFolderBO {
	d := dbClient.Model(&model.ProviderFolderBO{})
	d.Preload("Account")
	d.Where("deleted = 0 ")
	var list []model.ProviderFolderBO
	d.Find(&list)
	return list
}

func (p Provider) SelectListMeta() []model.ProviderFolderMeta {
	var rets []model.ProviderFolderMeta
	for _, t := range p.SelectList() {
		rets = append(rets, model.ProviderFolderMeta{ProviderFolderBO: t})
	}
	return rets
}
