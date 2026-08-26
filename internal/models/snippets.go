package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Insert 插入片段
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	const query = `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	reuslt, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := reuslt.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get 通过id查询片段,查询单条关键点使用queryrow不需要关闭连接
func (m *SnippetModel) Get(id int64) (*Snippet, error) {
	const query = `SELECT id, title, content, created, expires
	FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := new(Snippet)
	// 查询单条使用queryrow自动管理连接
	err := m.DB.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

// Latest 查询最后10个片段,查询多条使用query必须关闭连接
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	const query = `SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > UTC_TIMESTAMP()
	ORDER BY id DESC LIMIT 10`

	snippets := make([]*Snippet, 0)

	// 查询多条使用query，必须关闭连接
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := &Snippet{}
		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return snippets, nil

}
