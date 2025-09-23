package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var Err error
var mysqlLogger logger.Interface = logger.Default.LogMode(logger.Info)

func init() {
	username := "root"
	password := "xxh2023gkpku"
	host := "127.0.0.1"
	port := 3306
	Dbname := "gorm"
	timeout := "10s"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, //单数表名
			NoLowerCase:   false, //关闭大小写转换
		}, // Logger: mysqlLogger,
	})
	if err != nil {
		panic("连接失败，" + err.Error())
	}
	fmt.Println(db)
	DB = db
	Err = err
}
