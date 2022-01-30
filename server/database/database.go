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

var DB map[string]*gorm.DB

func Init() {
	DB = make(map[string]*gorm.DB, 2)
	db, err := gorm.Open(mysql.Open(getKR01DSN()), &gorm.Config{
		Logger: getGormLogger(),
		//NamingStrategy: schema.NamingStrategy{
		//}
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
	})
	if err != nil {
		panic(err)
	}
	DB["KR01"] = db

	db, err = gorm.Open(mysql.Open(getCN02DSN()), &gorm.Config{
		Logger: getGormLogger(),
		//NamingStrategy: schema.NamingStrategy{
		//}
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
	})
	if err != nil {
		panic(err)
	}
	DB["CN02"] = db
}

func getGormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Able color
		},
	)
}

func getKR01DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		server.DatabaseConfig.KR01User,
		server.DatabaseConfig.KR01Password,
		server.DatabaseConfig.KR01Host,
		server.DatabaseConfig.KR01Port,
		server.DatabaseConfig.KR01DatabaseName,
	)
}

func getCN02DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		server.DatabaseConfig.CN02User,
		server.DatabaseConfig.CN02Password,
		server.DatabaseConfig.CN02Host,
		server.DatabaseConfig.CN02Port,
		server.DatabaseConfig.CN02DatabaseName,
	)
}

//type Model struct {
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt gorm.DeletedAt
//}

// TODO : Gt, Gte, Lt, Lte, IfNotEqual

func Contains(db *gorm.DB, column string, value string) *gorm.DB {
	return db.Where(fmt.Sprintf("%s = ?", column), "%"+value+"%")
}

func NotContains(db *gorm.DB, column string, value string) *gorm.DB {
	return db.Not(fmt.Sprintf("%s = ?", column), "%"+value+"%")
}

func IfContains(db *gorm.DB, column string, value string) *gorm.DB {
	if value != "" {
		return Contains(db, column, value)
	}
	return db
}

func Equal(db *gorm.DB, column string, value interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s = ?", column), value)
}

func NotEqual(db *gorm.DB, column string, value interface{}) *gorm.DB {
	return db.Not(fmt.Sprintf("%s = ?", column), value)
}

func IfEqual(db *gorm.DB, column string, value interface{}) *gorm.DB {
	if value != "" && value != nil {
		return Equal(db, column, value)
	}
	return db
}
