package bootstrap

import (
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path/filepath"
)

func bootLog() {
	var logWriter io.Writer = os.Stdout
	logFile := viper.GetString("log.file")
	if logFile != "" {
		logFile, err := filepath.Abs(logFile)
		exception.Check(err)
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			_ = os.MkdirAll(filepath.Dir(logFile), 0777)
		}
		f, err := os.Create(logFile)
		exception.Check(err)
		logWriter = io.MultiWriter(f)
	}
	logger := log.Default()
	logger.SetPrefix("[SYS] ")
	logger.SetOutput(logWriter)
}
