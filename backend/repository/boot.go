package repository

import (
	"database/sql"
	"github.com/czy21/ndisk/exception"
	"github.com/czy21/ndisk/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var dbClient *gorm.DB

func Boot() {
	dbLoggerConfig := logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	}
	if viper.GetString("log.level") == "debug" {
		dbLoggerConfig.LogLevel = logger.Info
	}
	dbLogger := logger.New(log.StandardLogger(), dbLoggerConfig)
	dbConnect, err := sql.Open(viper.GetString("db.driver-name"), viper.GetString("db.url"))
	dbConnect.SetMaxIdleConns(5)
	dbConnect.SetMaxOpenConns(10)
	err = dbConnect.Ping()
	exception.Check(err)
	dbClient, err = gorm.Open(mysql.New(mysql.Config{
		Conn: dbConnect,
	}), &gorm.Config{
		Logger: dbLogger,
	})
}

type Base[T any] struct {
	*gorm.DB
}

func (b Base[T]) SelectPage(pageIndex int, pageSize int) ([]T, model.PageModel) {
	var list []T
	page := model.PageModel{PageIndex: pageIndex, PageSize: pageSize}
	b.Count(&page.Total)
	b.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&list)
	return list, page
}
