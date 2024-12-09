package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 接続プールの設定
	connStr := "postgres://postgres:postgres@localhost:5432/postgres"
	config, _ := pgxpool.ParseConfig(connStr)
	config.MaxConns = 10
	config.MinConns = 2

	// 接続プールを作成
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// 接続プールの状態をモニタリング
	go monitorPool(pool)

	// サンプルクエリを実行
	for i := 0; i < 5; i++ {
		go func() {
			conn, err := pool.Acquire(context.Background())
			if err != nil {
				log.Printf("Failed to acquire connection: %v", err)
				return
			}
			defer conn.Release()

			var now time.Time
			// err = conn.QueryRow(context.Background(), "SELECT NOW()").Scan(&now)
			log.Println("Executing pg_sleep(3)...")
			_, err = conn.Exec(context.Background(), "SELECT pg_sleep(3);")
			if err != nil {
				log.Printf("Query failed: %v", err)
				return
			}

			fmt.Printf("Current time: %v\n", now)
		}()
	}

	// しばらく実行
	time.Sleep(10 * time.Second)
}

// monitorPool は接続プールの状態を定期的にログに出力する
func monitorPool(pool *pgxpool.Pool) {
	for {
		stats := pool.Stat()
		fmt.Printf("Pool Stats - TotalConns: %d, IdleConns: %d, AcquiredConns: %d, MaxConns: %d\n",
			stats.TotalConns(),
			stats.IdleConns(),
			stats.AcquiredConns(),
			stats.MaxConns())
		time.Sleep(2 * time.Second)
	}
}
