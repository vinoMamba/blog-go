package database

import (
	"database/sql"
	"time"

	"github.com/vinoMamba/goblog/pkg/logger"
)

var DB *sql.DB

func Initialize() {
	initDB()
	createTables()
}

func initDB() {
	var err error
	connStr := "host=localhost user=mangosteen password=123456 dbname=mangosteen_dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// 准备数据库连接池
	DB, err = sql.Open("postgres", connStr)
	logger.LogError(err)

	// 设置最大连接数
	DB.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	DB.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = DB.Ping()
	logger.LogError(err)
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		body TEXT NOT NULL
		);`
	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
