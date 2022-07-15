package bootstrap

import (
	"github.com/czy21/cloud-disk-sync/cache"
	"github.com/czy21/cloud-disk-sync/controller"
	"github.com/czy21/cloud-disk-sync/repository"
)

func Boot() {
	bootConfig()
	bootLog()
	repository.Boot()
	cache.Boot()
	controller.Boot()
}
