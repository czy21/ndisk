package bootstrap

import (
	"github.com/czy21/ndisk/exception"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
)

func bootLog() {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true, DisableLevelTruncation: true, PadLevelText: true, ForceColors: true})
	confLevel := viper.GetString("log.level")
	confFile := viper.GetString("log.file")
	if confLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
	var output io.Writer = os.Stdout
	log.SetOutput(output)
	if confFile != "" {
		logFile, err := filepath.Abs(confFile)
		exception.Check(err)
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			_ = os.MkdirAll(filepath.Dir(logFile), 0777)
		}
		f, err := os.Create(logFile)
		exception.Check(err)
		log.SetOutput(io.MultiWriter(f))
		log.SetFormatter(&log.JSONFormatter{})
	}
}
