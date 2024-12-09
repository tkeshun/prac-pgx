package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres"
	config, err := pgxpool.ParseConfig(connStr)
	config.MaxConns = 10                        // 最大接続数、初期設定は4 もしくは runtime.NumCPU()の数、プールしてる接続が必要になると新しい接続は待機状態になる
	config.MinConns = 2                         // 最小接続数、設定数が低いとMaxConnLifetimeが過ぎるまでもしくは、ヘルスチェックがあるまで、Poolが空の状態になり接続が待機される可能性がある
	config.MaxConnIdleTime = 10 * time.Minute   // 各接続が維持される最大時間、この時間を超えると新しい接続に置き換えられる
	config.MaxConnIdleTime = 10 * time.Minute   // アイドル状態の接続が維持される最大時間、これを超えるとアイドル接続はクローズされる。
	config.HealthCheckPeriod = 15 * time.Second // 接続プールのヘルスチェック間隔、無効な接続を取り除いたり、MinConnsに足りない場合、接続を補充したりする

	// 各接続に関連する低レベル設定 pgconn.Config型
	config.ConnConfig.Host = "localhost"
	config.ConnConfig.Port = 5432
	config.ConnConfig.User = "postgres"
	config.ConnConfig.Password = "password"
	config.ConnConfig.Database = "DB"
	config.ConnConfig.ConnectTimeout = 5 * time.Second   // 接続待ち時間
	config.ConnConfig.RuntimeParams = map[string]string{ // 接続時に PostgreSQL に渡すランタイムパラメータ
		"search_path": "myschema",
		"timezone":    "UTC",
	}

	// Config型の定義
	// github.com/jackc/pgx/v5@v5.7.1/pgxpool/pool.go
	// // Config is the configuration struct for creating a pool. It must be created by [ParseConfig] and then it can be
	// // modified.
	// type Config struct {
	// 	ConnConfig *pgx.ConnConfig

	// 	// BeforeConnect is called before a new connection is made. It is passed a copy of the underlying pgx.ConnConfig and
	// 	// will not impact any existing open connections.
	// 	BeforeConnect func(context.Context, *pgx.ConnConfig) error

	// 	// AfterConnect is called after a connection is established, but before it is added to the pool.
	// 	AfterConnect func(context.Context, *pgx.Conn) error

	// 	// BeforeAcquire is called before a connection is acquired from the pool. It must return true to allow the
	// 	// acquisition or false to indicate that the connection should be destroyed and a different connection should be
	// 	// acquired.
	// 	BeforeAcquire func(context.Context, *pgx.Conn) bool

	// 	// AfterRelease is called after a connection is released, but before it is returned to the pool. It must return true to
	// 	// return the connection to the pool or false to destroy the connection.
	// 	AfterRelease func(*pgx.Conn) bool

	// 	// BeforeClose is called right before a connection is closed and removed from the pool.
	// 	BeforeClose func(*pgx.Conn)

	// 	// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	// 	MaxConnLifetime time.Duration

	// 	// MaxConnLifetimeJitter is the duration after MaxConnLifetime to randomly decide to close a connection.
	// 	// This helps prevent all connections from being closed at the exact same time, starving the pool.
	// 	MaxConnLifetimeJitter time.Duration

	// 	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	// 	MaxConnIdleTime time.Duration

	// 	// MaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU().
	// 	MaxConns int32

	// 	// MinConns is the minimum size of the pool. After connection closes, the pool might dip below MinConns. A low
	// 	// number of MinConns might mean the pool is empty after MaxConnLifetime until the health check has a chance
	// 	// to create new connections.
	// 	MinConns int32

	// 	// HealthCheckPeriod is the duration between checks of the health of idle connections.
	// 	HealthCheckPeriod time.Duration

	// 	createdByParseConfig bool // Used to enforce created by ParseConfig rule.
	// }

	if err != nil {
		log.Fatalln("Err")
	}

	printConfig(config)
}

func printConfig(config *pgxpool.Config) {
	fmt.Println("pgxpool.Config Settings:")
	fmt.Printf("  MaxConns: %d\n", config.MaxConns)
	fmt.Printf("  MinConns: %d\n", config.MinConns)
	fmt.Printf("  MaxConnLifetime: %v\n", config.MaxConnLifetime)
	fmt.Printf("  MaxConnIdleTime: %v\n", config.MaxConnIdleTime)
	fmt.Printf("  HealthCheckPeriod: %v\n", config.HealthCheckPeriod)

	fmt.Println("\npgconn.Config Settings:")
	fmt.Printf("  Host: %s\n", config.ConnConfig.Host)
	fmt.Printf("  Port: %d\n", config.ConnConfig.Port)
	fmt.Printf("  User: %s\n", config.ConnConfig.User)
	fmt.Printf("  Password: %s\n", config.ConnConfig.Password)
	fmt.Printf("  Database: %s\n", config.ConnConfig.Database)
	fmt.Printf("  ConnectTimeout: %v\n", config.ConnConfig.ConnectTimeout)
	fmt.Println("  RuntimeParams:")
	for key, value := range config.ConnConfig.RuntimeParams {
		fmt.Printf("    %s: %s\n", key, value)
	}
}
