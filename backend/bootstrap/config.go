package bootstrap

import (
	"github.com/czy21/ndisk/exception"
	"github.com/spf13/viper"
	"os"
)

func bootConfig() {
	appConf := os.Getenv("CONFIG_FILE")
	if appConf == "" {
		viper.SetConfigFile("app.yaml")
	} else {
		viper.SetConfigFile(appConf)
	}
	err := viper.ReadInConfig()
	exception.Check(err)
	viper.WatchConfig()
}
