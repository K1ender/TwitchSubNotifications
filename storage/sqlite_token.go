package storage

import (
	"database/sql"
	"subalertor/types"
	"time"
)

type TokenStore interface {
	AddTokens(user_id string, access_token string, refresh_token string, expires_at time.Time) error
	GetTokens(user_id string) (types.UserAccessToken, string, error)
	SetTokens(user_id string, access_token string, refresh_token string, expires_at time.Time) error
}

type SQLiteTokenStore struct {
	db *sql.DB
}

func NewSQLiteTokenStore(db *sql.DB) *SQLiteTokenStore {
	return &SQLiteTokenStore{
		db: db,
	}
}

func (s *SQLiteTokenStore) AddTokens(user_id string, access_token string, refresh_token string, expires_at time.Time) error {
	query := "INSERT INTO tokens (user_id, access_token, refresh_token, expires_at) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, user_id, access_token, refresh_token, time.Now().Add(time.Duration(expires_at.Unix())*time.Second).Unix())
	return err
}

func (s *SQLiteTokenStore) GetTokens(user_id string) (types.UserAccessToken, string, error) {
	query := "SELECT access_token, refresh_token FROM tokens WHERE user_id = ?"
	var access_token types.UserAccessToken
	var refresh_token string
	err := s.db.QueryRow(query, user_id).Scan(&access_token, &refresh_token)
	return access_token, refresh_token, err
}

func (s *SQLiteTokenStore) SetTokens(user_id string, access_token string, refresh_token string, expires_at time.Time) error {
	query := "UPDATE tokens SET access_token = ?, refresh_token = ?, expires_at = ? WHERE user_id = ?"
	_, err := s.db.Exec(query, access_token, refresh_token, time.Now().Add(time.Duration(expires_at.Unix())*time.Second).Unix(), user_id)
	return err
}
