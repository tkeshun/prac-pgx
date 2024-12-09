-- データベースのセットアップ

-- users テーブル: ユーザー情報を格納
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,         -- 自動インクリメントの主キー
    name VARCHAR(100) NOT NULL,    -- ユーザー名
    age INT NOT NULL,              -- 年齢
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 作成日時
);

-- json_table テーブル: JSON データを格納
CREATE TABLE IF NOT EXISTS json_table (
    id SERIAL PRIMARY KEY,         -- 自動インクリメントの主キー
    data JSONB NOT NULL,           -- JSON データ
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 作成日時
);

-- テストデータの挿入
INSERT INTO users (name, age) VALUES
    ('Alice', 25),
    ('Bob', 30),
    ('Charlie', 35)
ON CONFLICT DO NOTHING;

INSERT INTO json_table (data) VALUES
    ('{"key1": "value1", "key2": "value2"}'),
    ('{"key1": "value3", "key2": "value4"}')
ON CONFLICT DO NOTHING;