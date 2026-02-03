package ioc

import (
	"fmt"

	"moon/internal/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}

	c := Config{
		DSN: "root:root@tcp(localhost:3306)/mysql",
	}

	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败，原因 %v", err))
	}

	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("初始化数据库失败，原因 %v", err))
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(fmt.Errorf("初始化数据库表失败，原因 %v", err))
	}
	return db
}
