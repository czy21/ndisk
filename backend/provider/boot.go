package provider

import (
	_189 "github.com/czy21/cloud-disk-sync/provider/189"
)

var Providers map[string]FileSystem

func init() {
	Providers = make(map[string]FileSystem)
	Providers[string(Cloud189)] = _189.FileSystem{}
}
