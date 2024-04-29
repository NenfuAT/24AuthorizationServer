package model

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db = SqlConnect()

func SqlConnect() (database *gorm.DB) {
	var db *gorm.DB
	var err error

	// 環境変数を読み込む
	host := os.Getenv("MYSQL_CONTAINER_HOST")
	rootPassword := os.Getenv("MYSQL_ROOT_PASSWORD")
	databasename := os.Getenv("MYSQL_DATABASE")
	rootUser := os.Getenv("MYSQL_ROOT_USER")

	fmt.Println(host, rootPassword, databasename, rootUser)

	// MySQLのDSNフォーマットを使用
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		rootUser, rootPassword, host, databasename,
	)

	dialector := mysql.Open(dsn) // MySQL用のダイアレクタ

	fmt.Println(dsn) // DSNの表示
	if db, err = gorm.Open(dialector); err != nil {
		db = connect(dialector, 10) // リトライ関数
	}
	fmt.Println("db connected!!")

	return db // データベース接続を返す
}

func connect(dialector gorm.Dialector, count uint) *gorm.DB {
	var err error
	var db *gorm.DB
	if db, err = gorm.Open(dialector); err != nil {
		if count > 1 {
			time.Sleep(time.Second * 2)
			count--
			fmt.Printf("retry... count:%v\n", count)
			return connect(dialector, count) // 再帰的にリトライ
		}
		panic(err.Error()) // エラーを出力してパニック
	}
	return db
}
