package storage

import "database/sql"

type SQLiteFollowersStore struct {
	db *sql.DB
}

type FollowerModel struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	FollowedAt  int    `json:"followed_at"`
}

type FollowerStore interface {
	GetFollowers(userID string, offset, limit int) ([]FollowerModel, error)
	AddFollower(userID string, follower FollowerModel) error
	DeleteFollower(userID string, followerID string) error
}

func NewSQLiteFollowersStore(db *sql.DB) *SQLiteFollowersStore {
	return &SQLiteFollowersStore{
		db: db,
	}
}

func (s *SQLiteFollowersStore) GetFollowers(userID string, offset, limit int) ([]FollowerModel, error) {
	query := "SELECT id, display_name, username, followed_at FROM followers WHERE followed_to = ? LIMIT ? OFFSET ?"
	var followers []FollowerModel

	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var follower FollowerModel
		err := rows.Scan(&follower.ID, &follower.DisplayName, &follower.Username, &follower.FollowedAt)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(followers) == 0 {
		return []FollowerModel{}, nil
	}
	return followers, nil
}

func (s *SQLiteFollowersStore) AddFollower(userID string, follower FollowerModel) error {
	query := "INSERT INTO followers (display_name, username, followed_at, followed_to) VALUES ( ?, ?, ?, ?)"
	_, err := s.db.Exec(query, follower.DisplayName, follower.Username, follower.FollowedAt, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteFollowersStore) DeleteFollower(followerUsername, streamerID string) error {
	query := "DELETE FROM followers WHERE username = ? AND followed_to = ?"
	_, err := s.db.Exec(query, followerUsername, streamerID)
	return err
}
