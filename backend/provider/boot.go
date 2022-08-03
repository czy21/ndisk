package provider

import (
	"github.com/czy21/ndisk/model"
	_189 "github.com/czy21/ndisk/provider/189"
	"github.com/czy21/ndisk/provider/baidu"
)

var providers map[string]model.FileSystem

func Boot() {
	providers = make(map[string]model.FileSystem)
	providers[string(Cloud189)] = _189.FileSystem{}
	providers[string(CloudBaiDu)] = baidu.FileSystem{}
}

func GetProviders() map[string]model.FileSystem {
	return providers
}
