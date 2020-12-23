package model

import (
	"goblog/pkg/logger"
	"goblog/pkg/types"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// BaseModel 模型基类
type BaseModel struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"`
	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

// GetStringID 获取string类型的ID
func (model BaseModel) GetStringID() string {
	return types.Uint64ToString(model.ID)
}

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
