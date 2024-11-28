package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 接続文字列（URI形式）
	connStr := "postgres://postgres:postgres@localhost:5432/postgres"

	// 接続プールを作成
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	// 接続が成功したことを確認
	log.Println("Successfully connected to the database using pgxpool!")
}
