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

	// INSERT
	insert, err := pool.Exec(context.Background(), "INSERT INTO users (name, age) VALUES ($1, $2)", "Alice", 25)
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	log.Printf("Rows affected: %d", insert.RowsAffected()) // クエリを実行した際に影響を受けた（変更された）行数を返す、pgx.CommandTag 型のメソッド

	// INSERT
	insertID1, err := pool.Exec(context.Background(), "INSERT INTO users (id,name, age) VALUES (1,$1, $2)", "Alice", 25)
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	log.Printf("Rows affected: %d", insertID1.RowsAffected()) // クエリを実行した際に影響を受けた（変更された）行数を返す、pgx.CommandTag 型のメソッド

	// SELECT(単一行)
	var name string
	var age int
	errQueryRow := pool.QueryRow(context.Background(), "SELECT name, age FROM users WHERE id=$1", 1).Scan(&name, &age)
	if errQueryRow != nil {
		log.Fatalf("QueryRow failed: %v", err)
	}
	log.Printf("Name: %s, Age: %d", name, age)

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

	// UPDATE
	update, err := pool.Exec(context.Background(), "UPDATE users SET age=$1 WHERE name=$2", 26, "Alice")
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Rows affected: %d", update.RowsAffected())

	// DELETE
	delete, err := pool.Exec(context.Background(), "DELETE FROM users WHERE name=$1", "Alice")
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Rows affected: %d", delete.RowsAffected())

	// pgx.CommandTag型
	// 	Select: 実行したクエリが SELECT かどうかを判定
	if rows.CommandTag().Select() {
		log.Println("rowsはSELECT Query!!")
	}
	if !rows.CommandTag().Insert() {
		log.Println("rowsはINSERTクエリではない")
	}
	// Insert: 実行したクエリが INSERT かどうかを判定
	if insert.Insert() {
		log.Println("INSERT")
	}
	// Update: 実行したクエリが UPDATE かどうかを判定
	if update.Update() {
		log.Println("UPDATE")
	}
	// Delete: 実行したクエリが DELETE かどうかを判定
	if delete.Delete() {
		log.Println("DELETE")
	}
}
