package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewDBShard() *ShardManager {
	shards := make(map[int]*sql.DB)

	proxySql0 := NewDB(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("PROXY_SQL_0_HOST"), os.Getenv("PROXY_SQL_0_PORT"), os.Getenv("DB_NAME"))
	proxySql1 := NewDB(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("PROXY_SQL_1_HOST"), os.Getenv("PROXY_SQL_1_PORT"), os.Getenv("DB_NAME"))

	shards[0] = proxySql0
	shards[1] = proxySql1

	return &ShardManager{
		Shards: shards,
	}
}

func NewDB(user, password, host, port, dbName string) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user,
		password,
		host,
		port,
		dbName,
	)
	log.Println("Shard Connection url=", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 🔥 Pooling (critical)
	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(10)
	// db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Shard not reachable:", err)
	}

	log.Printf("✅ Shard Connected, URL=%s", dsn)
	return db
}
