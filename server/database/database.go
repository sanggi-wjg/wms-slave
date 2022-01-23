package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"wms_slave/server"
)

var DB *gorm.DB

func Init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Able color
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(getDSN()), &gorm.Config{
		Logger: newLogger,
		//NamingStrategy: schema.NamingStrategy{
		//
		//}
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
	})
	if err != nil {
		panic(err)
	}
}

func getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		server.DatabaseConfig.User,
		server.DatabaseConfig.Password,
		server.DatabaseConfig.Host,
		server.DatabaseConfig.Port,
		server.DatabaseConfig.DatabaseName,
	)
}

//type Model struct {
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt gorm.DeletedAt
//}
