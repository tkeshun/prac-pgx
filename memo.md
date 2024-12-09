# pgxの基礎

https://pkg.go.dev/github.com/jackc/pgx/v5

## pgxとは

Postgresql用のSQLドライバー
Postgresql固有の機能へアクセスするための低レベルのインターフェース

## 環境構築


### goプロジェクト
- project初期化

`go mod init basic-pgx`

- go.modへ追記

`go get github.com/jackc/pgx/v5`

```
go get github.com/jackc/pgx/v5
go: added github.com/jackc/pgpassfile v1.0.0
go: added github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761
go: added github.com/jackc/pgx/v5 v5.7.1
go: added golang.org/x/crypto v0.27.0
go: added golang.org/x/text v0.18.0
```


### DB

compose.yamlを参照  
起動コマンド  
`docker compose up -d`


## DB接続

`postgres://postgres:postgres@localhost:5432/postgres`

- postgres://

    PostgreSQL用の接続URIスキーム  
    ※接頭辞について  
        postgreSQLデータベースサーバーに接続する際、以下の2つの接頭辞が使われる  
        - postgres://
        - postgresql://    
        クライアントライブラリがPostgresqlへの接続であることを解釈してつなぐ
    postgreSQLの通信プロトコルはアプリケーション層で動作する
    ネットワークの低レベルでは、通常 TCP（Transmission Control Protocol）が使われる

    postgres://がついてるか判断するロジックが`/home/shun/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.1/pgconn/config.go > ParseConfigWithOptions` にある

    ```
    func ParseConfigWithOptions(connString string, options ParseConfigOptions) (*Config, error) {
	defaultSettings := defaultSettings()
	envSettings := parseEnvSettings()

	connStringSettings := make(map[string]string)
	if connString != "" { //  文字列（URI）が指定されてた場合
		var err error
		// connString may be a database URL or in PostgreSQL keyword/value format
		if strings.HasPrefix(connString, "postgres://") || strings.HasPrefix(connString, "postgresql://") { // postgres://もしくはpostgresql://が含まれてたらURI形式だと判断する
			connStringSettings, err = parseURLSettings(connString) // URI用のパース関数
			if err != nil {
				return nil, &ParseConfigError{ConnString: connString, msg: "failed to parse as URL", err: err} // うまくパースできない場合、URLが入ってきてると判断したブロックなのでURLがだめだと返す
			}
		} else {
            // 文字列指定、postgresなしの場合、key/value形式（例：host=localhost port=5432 user=postgres password=postgres dbname=postgres）だと判断する
			connStringSettings, err = parseKeywordValueSettings(connString) // key-valueのパース用関数  
			if err != nil {
				return nil, &ParseConfigError{ConnString: connString, msg: "failed to parse as keyword/value", err: err} // だめだった場合、正しいkey-valueではないので、その旨+エラーで返す
			}
		}
	}
    ```

    ""を指定した場合は環境変数が使われる

    pgx.ConnConfigを使いたい場合は、`pgx.ConnectConfig`を使う  
    ```
    	connConfig := &pgconn.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		ConnectTimeout: 5 * time.Second, // カスタムタイムアウト
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
    ```

- postgres:postgres

      ユーザー名:パスワード

- @localhost

    ホスト名  
    他のサーバーを指定する場合はIPアドレスやドメイン名を置き換える

- :5432

    データベースサーバーがリッスンしてるポート番号
    ポート番号5432は、Postgresqlのプロトコルをサポートするサーバ用のTCPポートとしてIANAに登録されている  
    `https://www.postgresql.jp/document/16/html/protocol.html`  



- /postgres

    接続先のデータベース名  
    接続したいデータベース名を指定する



### 同時実行性安全な接続プール

*pgx.Connは同時実行が安全ではないため、github.com/jackc/pgx/v5/pgxpoolを使う必要あり

Pool型の値を返す  
Pool.pがリソースプール（puddle.Pool型）
Pool.muxを持つ  

リソースの取得 (Acquire)

acquireSem を使ってリソースの取得許可を確認  
未使用リソースがある場合は idleResources から取得  
リソースが不足している場合は Constructor を使って新しいリソースを生成  
新規接続時には接続リソースを作り、allResourcesにリソースを追加する  
/home/shun/go/pkg/mod/github.com/jackc/puddle/v2@v2.2.2/pool.go  

```
func (p *Pool[T]) createNewResource() *Resource[T] {
	res := &Resource[T]{
		pool:           p,
		creationTime:   time.Now(),
		lastUsedNano:   nanotime(),
		poolResetCount: p.resetCount,
		status:         resourceStatusConstructing,
	}

	p.allResources.append(res)
	p.destructWG.Add(1)

	return res
}
```

サンプル

```
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
```


### プール設定

```
	config.MaxConns = 10                        // 最大接続数、初期設定は4 もしくは runtime.NumCPU()の数、プールしてる接続が必要になると新しい接続は待機状態になる
	config.MinConns = 2                         // 最小接続数、設定数が低いとMaxConnLifetimeが過ぎるまでもしくは、ヘルスチェックがあるまで、Poolが空の状態になり接続が待機される可能性がある
	config.MaxConnIdleTime = 10 * time.Minute   // 各接続が維持される最大時間、この時間を超えると新しい接続に置き換えられる
	config.MaxConnIdleTime = 10 * time.Minute   // アイドル状態の接続が維持される最大時間、これを超えるとアイドル接続はクローズされる。
	config.HealthCheckPeriod = 15 * time.Second // 接続プールのヘルスチェック間隔、無効な接続を取り除いたり、MinConnsに足りない場合、接続を補充したりする
```

### 接続プールのモニタリング

`pool.Stat()`で接続プールの統計情報を取得できる

StatはPoolに格納されてる変数から情報を抜いてる

（例）pgxpool-monitor/main.go

- TotalConns: プール内の総接続数  
- IdleConns: アイドル状態（未使用）の接続数  
- AcquiredConns: 使用中の接続数  
- MaxConns: プール内の最大接続数  

```
$ go run main.go 
Pool Stats - TotalConns: 7, IdleConns: 0, AcquiredConns: 0, MaxConns: 10
Current time: 2024-11-30 00:38:53.538471 +0900 JST
Current time: 2024-11-30 00:38:53.53847 +0900 JST
Current time: 2024-11-30 00:38:53.538464 +0900 JST
Current time: 2024-11-30 00:38:53.539802 +0900 JST
Current time: 2024-11-30 00:38:53.540358 +0900 JST
Pool Stats - TotalConns: 7, IdleConns: 7, AcquiredConns: 0, MaxConns: 10
Pool Stats - TotalConns: 7, IdleConns: 7, AcquiredConns: 0, MaxConns: 10
Pool Stats - TotalConns: 7, IdleConns: 7, AcquiredConns: 0, MaxConns: 10
Pool Stats - TotalConns: 7, IdleConns: 7, AcquiredConns: 0, MaxConns: 10
```

 接続中の統計情報を見る

```
 go run main.go 
Pool Stats - TotalConns: 6, IdleConns: 0, AcquiredConns: 0, MaxConns: 10
2024/11/30 00:42:50 Executing pg_sleep(3)...
2024/11/30 00:42:50 Executing pg_sleep(3)...
2024/11/30 00:42:50 Executing pg_sleep(3)...
2024/11/30 00:42:50 Executing pg_sleep(3)...
2024/11/30 00:42:50 Executing pg_sleep(3)...
Pool Stats - TotalConns: 7, IdleConns: 2, AcquiredConns: 5, MaxConns: 10
Current time: 0001-01-01 00:00:00 +0000 UTC
Current time: 0001-01-01 00:00:00 +0000 UTC
Current time: 0001-01-01 00:00:00 +0000 UTC
Current time: 0001-01-01 00:00:00 +0000 UTC
Current time: 0001-01-01 00:00:00 +0000 UTC
Pool Stats - TotalConns: 7, IdleConns: 7, AcquiredConns: 0, MaxConns: 10
Pool Stats - TotalConns: 7, IdleConns: 7, AcquiredConns: 0, MaxConns: 10
Pool Stats - TotalConns: 7, IdleConns: 7, AcquiredConns: 0, MaxConns: 10
```

## 基本的なクエリ操作

### QueryRow

1行限定、error型しか返さない
Scanで指定した変数に値を格納する

```
	// SELECT(単一行)
	var name string
	var age int
	errQueryRow := pool.QueryRow(context.Background(), "SELECT name, age FROM users WHERE id=$1", 1).Scan(&name, &age)
	if errQueryRow != nil {
		log.Fatalf("QueryRow failed: %v", err)
	}
	log.Printf("Name: %s, Age: %d", name, age)

```

### Query

複数行SELECTする
返り値の変数のnext()を使って値を取り出す。
```
	// SELECT（複数行）
	rows, err := pool.Query(context.Background(), "SELECT id, name, age FROM users")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatalf("Row scan failed: %v", err)
		}
		log.Printf("ID: %d, Name: %s, Age: %d", id, name, age)
	}
```

### Exec

### CommandTag型

CommandTag
pgx.Exec や pgx.CopyFrom の結果として返される  
CommandTag 型には、クエリの実行結果や影響を受けた行数に関する情報が含まれる  

主なメソッド
RowsAffected: 影響を受けた行数を返す
Select: 実行したクエリが SELECT かどうかを判定
Insert: 実行したクエリが INSERT かどうかを判定
Update: 実行したクエリが UPDATE かどうかを判定
Delete: 実行したクエリが DELETE かどうかを判定

例）
```
	// Delete: 実行したクエリが DELETE かどうかを判定
	if delete.Delete() {
		log.Println("DELETE")
	}
```

## トランザクション

### トランザクションの基本

- Begin, Commit

- savepoint

### トランザクションのエラーハンドリング


- Rollback

### 分離レベルの変更方法



## Batch API(バッチ処理)

### バッチの基本

- pgx.Batch を使った複数クエリの実行

- 結果をpgx.BatchResultsで取得する方法

### エラーハンドリング

- クエリ失敗時の処理
- 部分的な結果の取り扱い

### 実用例

- 大量データの挿入と更新
- バッチ処理でのトランザクション管理

## COPYプロトコル

### COPYプロトコルとは
PostgreSQLの高速データ操作機能。
大量データの効率的な読み書き。

### Pgxでの実装

- pgx.CopyFrom を使ったデータ挿入。
- カスタムリーダーを利用した動的データ処理

### パフォーマンス向上テクニック

- データ型とバッファサイズの最適化
- エラーハンドリングとリカバリ戦略


## ログとトレース

## 高パフォーマンス設計


### パフォーマンス最適化

- プリペアドステートメントの活用
- 適切な接続プールサイズの設定
- パラレルクエリの活用
- 並列処理でのPgxの使用例
- 複数クエリの同時実行設計
- クエリの実行時間計測
- 高頻度クエリの分析


### 関連ライブラリ

- 
    PostgreSQL ワイヤ プロトコルをモックするサーバーを作成する機能を提供
    pgproto3 と pgmock を組み合わせることで、PostgreSQL プロキシまたは MitM (カスタム接続プーラーなど) を実装するために必要な基本的なツールのほとんどが提供される
    https://github.com/jackc/pgmock

- 