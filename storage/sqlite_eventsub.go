package storage

import (
	"database/sql"
	"fmt"
	"strings"
)

type EventSubModel struct {
	ID        int            `json:"id"`
	Condition ConditionModel `json:"condition"`
	Type      string         `json:"type"`
}

func (e EventSubModel) PrettyPrint() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("\nID: " + fmt.Sprint(e.ID) + "\n")
	stringBuilder.WriteString("Type: " + e.Type + "\n")
	if e.Condition.BroadcasterID != nil {
		stringBuilder.WriteString("BroadcasterID: " + *e.Condition.BroadcasterID + "\n")
	}
	if e.Condition.UserID != nil {
		stringBuilder.WriteString("UserID: " + *e.Condition.UserID + "\n")
	}
	if e.Condition.BroadCastUserID != nil {
		stringBuilder.WriteString("BroadCastUserID: " + *e.Condition.BroadCastUserID + "\n")
	}
	if e.Condition.ModeratorUserID != nil {
		stringBuilder.WriteString("ModeratorUserID: " + *e.Condition.ModeratorUserID + "\n")
	}
	return stringBuilder.String()
}

type ConditionModel struct {
	ID              int     `json:"-"`
	BroadCastUserID *string `json:"broadcaster_user_id,omitempty"`
	ModeratorUserID *string `json:"moderator_user_id,omitempty"`
	BroadcasterID   *string `json:"broadcaster_id,omitempty"`
	UserID          *string `json:"user_id,omitempty"`
}

type EventSubStore interface {
	AddEventSubscription(eventSub EventSubModel) error
	GetAllEventSubscriptions() ([]EventSubModel, error)
	GetSubscribedEvents(userID string) ([]EventSubModel, error)
}

type SQLiteEventSubStore struct {
	db *sql.DB
}

func NewSQLiteEventSubStore(db *sql.DB) *SQLiteEventSubStore {
	return &SQLiteEventSubStore{
		db: db,
	}
}

func (s *SQLiteEventSubStore) AddEventSubscription(eventSub EventSubModel) error {
	query := "INSERT INTO conditions (broadcaster_id, user_id, broadcast_user_id, moderator_user_id) VALUES (?, ?, ?, ?) RETURNING id"

	var id int
	err := s.db.QueryRow(query, eventSub.Condition.BroadcasterID, eventSub.Condition.UserID, eventSub.Condition.BroadCastUserID, eventSub.Condition.ModeratorUserID).Scan(&id)

	if err != nil {
		return err
	}

	query = "INSERT INTO events (type, condition_id) VALUES (?, ?) RETURNING id"

	var eventID int
	err = s.db.QueryRow(query, eventSub.Type, id).Scan(&eventID)

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteEventSubStore) GetAllEventSubscriptions() ([]EventSubModel, error) {
	query := "SELECT e.id, e.type, c.broadcaster_id, c.user_id, c.broadcast_user_id, c.moderator_user_id FROM events e JOIN conditions c ON e.condition_id = c.id"
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	var eventSubs []EventSubModel
	for rows.Next() {
		var eventSub EventSubModel
		err := rows.Scan(&eventSub.ID, &eventSub.Type, &eventSub.Condition.BroadcasterID, &eventSub.Condition.UserID, &eventSub.Condition.BroadCastUserID, &eventSub.Condition.ModeratorUserID)
		if err != nil {
			return nil, err
		}
		eventSubs = append(eventSubs, eventSub)
	}

	return eventSubs, nil
}

func (s *SQLiteEventSubStore) GetSubscribedEvents(userID string) ([]EventSubModel, error) {
	query := "SELECT e.id, e.type, c.broadcaster_id, c.user_id, c.broadcast_user_id, c.moderator_user_id FROM events e JOIN conditions c ON e.condition_id = c.id WHERE c.user_id = ?"
	rows, err := s.db.Query(query, userID)

	if err != nil {
		return nil, err
	}

	var eventSubs []EventSubModel
	for rows.Next() {
		var eventSub EventSubModel
		err := rows.Scan(&eventSub.ID, &eventSub.Type, &eventSub.Condition.BroadcasterID, &eventSub.Condition.UserID, &eventSub.Condition.BroadCastUserID, &eventSub.Condition.ModeratorUserID)
		if err != nil {
			return nil, err
		}
		eventSubs = append(eventSubs, eventSub)
	}

	return eventSubs, nil
}
