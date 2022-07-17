package bootstrap

import (
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/controller"
	"github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/repository"
)

func Boot() {
	bootConfig()
	bootLog()
	repository.Boot()
	cache.Boot()
	http.Boot()
	controller.Boot()
}
