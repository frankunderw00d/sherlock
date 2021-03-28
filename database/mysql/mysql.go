package mysql

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ()

const ()

var (
	defaultDB *gorm.DB
)

func init() {}

// 初始化 DB
func InitializeMySQL(user, password, host, port, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return err
	}

	defaultDB = db

	return nil
}

// 注册表
func RegisterTable(tables ...interface{}) error {
	if defaultDB == nil {
		return errors.New("default DB is nil,please initialize it first")
	}
	return defaultDB.Migrator().AutoMigrate(tables...)
}

// 获取 DB 操控
func DB() *gorm.DB {
	return defaultDB
}
