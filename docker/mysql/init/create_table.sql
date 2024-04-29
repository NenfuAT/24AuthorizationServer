-- データベースを作成
CREATE DATABASE IF NOT EXISTS auth_server;

-- データベースに切り替え
USE auth_server;

-- テーブルを作成
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255),
    password VARCHAR(50),
    name_ja VARCHAR(50),
    given_name VARCHAR(50),
    family_name VARCHAR(50),
    locale VARCHAR(50)
);
