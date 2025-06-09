package database

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Connect() {
	// 환경 변수에서 데이터베이스 연결 정보 가져오기
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 데이터베이스 연결
	conn, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	// 커넥션 풀 설정
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(time.Hour)
	conn.SetConnMaxIdleTime(10 * time.Minute)

	db = conn

	slog.Info("database connection established successfully")
}

func GetConnection() *sqlx.DB {
	if db == nil {
		slog.Info("database connection is not established")
	}
	return db
}
