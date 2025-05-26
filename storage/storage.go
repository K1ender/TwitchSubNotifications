package storage

import "database/sql"

type Storage struct {
	UserStore     UserStore
	SessionStore  SessionStore
	TokenStore    TokenStore
	EventSubStore EventSubStore
	FollowerStore FollowerStore
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStore:     NewSQLiteUserStore(db),
		SessionStore:  NewSQLiteSessionStore(db),
		TokenStore:    NewSQLiteTokenStore(db),
		EventSubStore: NewSQLiteEventSubStore(db),
		FollowerStore: NewSQLiteFollowersStore(db),
	}
}
