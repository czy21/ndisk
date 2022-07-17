package provider

import (
	_189 "github.com/czy21/ndisk/provider/189"
	"github.com/czy21/ndisk/provider/local"
)

var providers map[string]FileSystem

func Boot() {
	providers = make(map[string]FileSystem)
	providers[string(Cloud189)] = _189.FileSystem{}
	providers[string(CloudLocal)] = local.NewFS()
}

func GetProviders() map[string]FileSystem {
	return providers
}
