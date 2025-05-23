package storage

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

type UserStore interface {
	CreateUser(id, username string) (UserModel, error)
	FindUser(username string) (UserModel, error)
	FindUserByID(id string) (UserModel, error)
	DeleteUser(id string) error
}

type UserModel struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type SQLiteUserStore struct {
	db *sql.DB
}

func NewSQLiteUserStore(db *sql.DB) *SQLiteUserStore {
	return &SQLiteUserStore{
		db: db,
	}
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

func (s *SQLiteUserStore) CreateUser(id, username string) (UserModel, error) {
	query := "INSERT INTO users (id, username) VALUES (?, ?) RETURNING id, username"
	var user UserModel

	err := s.db.QueryRow(query, id, username).Scan(&user.ID, &user.Username)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return user, ErrUserAlreadyExists
			}
		}
		return user, err
	}

	return user, nil
}

func (s *SQLiteUserStore) FindUser(username string) (UserModel, error) {
	query := "SELECT id, username FROM users WHERE username = ?"
	var user UserModel

	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *SQLiteUserStore) FindUserByID(id string) (UserModel, error) {
	query := "SELECT id, username FROM users WHERE id = ?"
	var user UserModel

	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *SQLiteUserStore) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}
