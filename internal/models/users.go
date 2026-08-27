package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID      int
	Name    string
	Email   string
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	const query = `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(query, name, email, hash)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hash []byte
	const query = `SELECT id, hashed_password FROM users WHERE email = ?`
	// 根据用户邮箱查询用户密码
	err := m.DB.QueryRow(query, email).Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		}
		return 0, err
	}

	// 校验用户密码
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return 0, ErrInvalidCredentials
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	const query = `SELECT EXISTS(SELECT TRUE FROM users WHERE id = ?)`
	err := m.DB.QueryRow(query, id).Scan(&exists)
	return exists, err
}
