package bootstrap

import (
	"github.com/czy21/cloud-disk-sync/controller"
	"github.com/czy21/cloud-disk-sync/repository"
)

func Boot() {
	bootConfig()
	bootLog()
	repository.Boot()
	controller.Boot()
}
