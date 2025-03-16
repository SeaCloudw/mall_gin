package database

import (
	// 导入 models 包

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// 初始化 , 连接mysql , 并获取db
func InitDB() {
	dsn := "root:root@tcp(127.0.0.1:3306)/mall_gin?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("无法连接数据库")
	}

	// 自动迁移表结构
	// 自动迁移模型
	// if err := db.AutoMigrate(&models.Customer{}); err != nil {
	// 	panic("failed to migrate database")
	// }
	// db.AutoMigrate(&models.Customer{})
}

func GetDB() *gorm.DB {
	return db
}
