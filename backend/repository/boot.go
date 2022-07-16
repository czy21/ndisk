package repository

import (
	"database/sql"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/czy21/cloud-disk-sync/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var dbClient *gorm.DB

func Boot() {
	dbLogger := logger.New(log.StandardLogger(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	dbConnect, err := sql.Open(viper.GetString("db.driver-name"), viper.GetString("db.url"))
	dbConnect.SetMaxIdleConns(5)
	dbConnect.SetMaxOpenConns(10)
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
