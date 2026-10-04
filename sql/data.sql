-- 创建能完整保存 Unicode 的数据库。
CREATE DATABASE snippetbox
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE snippetbox;

CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

-- 使用 '%' 允许任意来源连接（Docker 端口映射进来的连接不属于 localhost）
CREATE USER 'web'@'%';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'%';
ALTER USER 'web'@'%' IDENTIFIED BY 'pass';

-- 插入数据
INSERT INTO snippets (title, content, created, expires) VALUES
('First snippet', 'Hello from the database.', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)),
('Seven day snippet', 'This record expires in seven days.', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY));

-- 会话状态
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

-- 用户表
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users
  ADD CONSTRAINT users_uc_email UNIQUE (email);