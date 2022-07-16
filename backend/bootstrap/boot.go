package bootstrap

import (
	"github.com/czy21/cloud-disk-sync/cache"
	"github.com/czy21/cloud-disk-sync/controller"
	"github.com/czy21/cloud-disk-sync/http"
	"github.com/czy21/cloud-disk-sync/repository"
)

func Boot() {
	bootConfig()
	bootLog()
	repository.Boot()
	cache.Boot()
	http.Boot()
	controller.Boot()
}
