package bootstrap

import (
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/model"
	"time"

	"gorm.io/gorm"
)

// SetupDB 数据库设置
func SetupDB() {
	db := model.ConnectDB()

	sqlDB, _ := db.DB()

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲数
	sqlDB.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	migration(db)
}

// migration 执行数据表语句
func migration(db *gorm.DB) {
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
	)
}
