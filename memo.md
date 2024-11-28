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

### 接続プールのモニタリング

## 基本的なクエリ操作

### QueryRow

### Query

## トランザクション

### トランザクションの基本

- Begin, Commit

- savepoint

### トランザクションのエラーハンドリング


- Rollback


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