package test

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type Keys struct {
	Host string
	Port int32
	User string
	Pass string
	DB   string
}

func TestPostgres(t *testing.T) {
	keys := Keys{
		Host: "192.168.2.182",
		Port: 5432,
		User: "postgres",
		Pass: "msdnmm",
		DB:   "postgres",
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", keys.Host, keys.User, keys.Pass, keys.DB, keys.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Ping数据库
	sqlDB, err := db.DB()
	if err != nil {
		panic("无法获取数据库连接")
	}
	err = sqlDB.Ping()
	if err != nil {
		panic("无法ping数据库")
	}
}
