package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDBCluster() *DBCluster {
	primary := NewDB(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	replica := NewDB(os.Getenv("REPLICA_DB_USER"), os.Getenv("REPLICA_DB_PASSWORD"), os.Getenv("REPLICA_DB_HOST"), os.Getenv("REPLICA_DB_PORT"), os.Getenv("REPLICA_DB_NAME"))
	return &DBCluster{
		Primary:  primary,
		Replicas: []*sql.DB{replica},
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
	log.Println("Database Connection url=", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 🔥 Pooling (critical)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Printf("✅ DB Connected, URL=%s", dsn)
	return db
}
