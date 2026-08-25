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

CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

-- 插入数据
INSERT INTO snippets (title, content, created, expires) VALUES
('First snippet', 'Hello from the database.', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)),
('Seven day snippet', 'This record expires in seven days.', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY));
