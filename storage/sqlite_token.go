package storage

import (
	"database/sql"
	"subalertor/types"
	"time"
)

type TokenStore interface {
	AddTokens(userID, accessToken, refreshToken string, expires_at time.Time) error
	GetTokens(userID string) (types.UserAccessToken, string, error)
	SetTokens(userID, accessToken, refreshToken string, expires_at time.Time) error
	UpdateAccessToken(refreshToken, accessToken string) error
}

type SQLiteTokenStore struct {
	db *sql.DB
}

func NewSQLiteTokenStore(db *sql.DB) *SQLiteTokenStore {
	return &SQLiteTokenStore{
		db: db,
	}
}

func (s *SQLiteTokenStore) AddTokens(userID string, accessToken string, refreshToken string, expires_at time.Time) error {
	query := "INSERT INTO tokens (user_id, access_token, refresh_token, expires_at) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, userID, accessToken, refreshToken, time.Now().Add(time.Duration(expires_at.Unix())*time.Second).Unix())
	return err
}

func (s *SQLiteTokenStore) GetTokens(userID string) (types.UserAccessToken, string, error) {
	query := "SELECT access_token, refresh_token FROM tokens WHERE user_id = ?"
	var accessToken types.UserAccessToken
	var refreshToken string
	err := s.db.QueryRow(query, userID).Scan(&accessToken, &refreshToken)
	return accessToken, refreshToken, err
}

func (s *SQLiteTokenStore) SetTokens(userID string, accessToken string, refreshToken string, expiresAt time.Time) error {
	query := "UPDATE tokens SET access_token = ?, refresh_token = ?, expires_at = ? WHERE user_id = ?"
	_, err := s.db.Exec(query, accessToken, refreshToken, time.Now().Add(time.Duration(expiresAt.Unix())*time.Second).Unix(), userID)
	return err
}

func (s *SQLiteTokenStore) UpdateAccessToken(refreshToken string, access_token string) error {
	query := "UPDATE tokens SET access_token = ? WHERE refresh_token = ?"
	_, err := s.db.Exec(query, access_token, refreshToken)
	return err
}
