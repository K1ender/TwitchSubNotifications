package storage

import (
	"database/sql"
	"time"
)

type SessionStore interface {
	CreateSession(userID, id string, expiresAt time.Time) error
	DeleteSession(id string) error
	DeleteAllSessions(userID string) error
	FindSession(id string) (SessionModel, error)
	ExtendSession(id string, expiresAt time.Time) (SessionModel, error)
}

type SessionModel struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ExpiresAt int    `json:"expires_at"`
}

type SQLiteSessionStore struct {
	db *sql.DB
}

func NewSQLiteSessionStore(db *sql.DB) *SQLiteSessionStore {
	return &SQLiteSessionStore{
		db: db,
	}
}

func (s *SQLiteSessionStore) CreateSession(userID, id string, expiresAt time.Time) error {
	query := "INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)"
	_, err := s.db.Exec(query, id, userID, expiresAt.Unix())
	return err
}

func (s *SQLiteSessionStore) DeleteSession(id string) error {
	query := "DELETE FROM sessions WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SQLiteSessionStore) DeleteAllSessions(userID string) error {
	query := "DELETE FROM sessions WHERE user_id = ?"
	_, err := s.db.Exec(query, userID)
	return err
}

func (s *SQLiteSessionStore) FindSession(id string) (SessionModel, error) {
	query := "SELECT id, user_id, expires_at FROM sessions WHERE id = ?"
	var session SessionModel

	err := s.db.QueryRow(query, id).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return session, err
	}

	return session, nil
}

func (s *SQLiteSessionStore) ExtendSession(id string, expiresAt time.Time) (SessionModel, error) {
	query := "UPDATE sessions SET expires_at = ? WHERE id = ? RETURNING id, user_id, expires_at"

	var session SessionModel
	err := s.db.QueryRow(query, expiresAt.Unix(), id).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return session, err
	}

	return session, nil
}
