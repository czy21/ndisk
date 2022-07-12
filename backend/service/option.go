package service

import "github.com/czy21/cloud-disk-sync/model"

type Option struct {
}

func (Option) FindByKeys(query model.OptionQuery) map[string]interface{} {
	optionRet := make(map[string]interface{})
	optionAll := Option{}.all()
	for _, t := range query.Keys {
		v := optionAll[t]
		if v == nil {
			v = []model.SimpleItemModel[any]{}
		}
		optionRet[t] = v
	}
	return optionRet
}

func (Option) all() map[string]interface{} {
	option := make(map[string]interface{})
	return option
}
