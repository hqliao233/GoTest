package model

import (
	"goblog/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB 数据库操作对象
var DB *gorm.DB

// ConnectDB 连接数据库
func ConnectDB() *gorm.DB {
	var err error

	// 数据库连接配置
	config := mysql.New(mysql.Config{
		DSN: "root:root@tcp(127.0.0.1:3306)/gotest?charset=utf8&parseTime=True&loc=Local",
	})

	// 获取数据库操作对象并使用日志记录
	DB, err = gorm.Open(config, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})

	logger.LogError(err)

	return DB
}
