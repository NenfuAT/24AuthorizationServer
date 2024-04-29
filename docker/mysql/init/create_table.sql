-- データベースを作成
CREATE DATABASE IF NOT EXISTS auth_server;

-- データベースに切り替え
USE auth_server;

-- テーブルを作成
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) PRIMARY KEY,
    email VARCHAR(255),
    password VARCHAR(50),
    nameJa VARCHAR(50),
    givenName VARCHAR(50),
    familyName VARCHAR(50),
    locale VARCHAR(50)
);
