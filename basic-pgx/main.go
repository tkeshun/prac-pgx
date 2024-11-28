package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	// プロトコル://ユーザー名:パスワード@ホスト名:port/DB名
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	log.Println("Connected to the database using environment variables!")
}
