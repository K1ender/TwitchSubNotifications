package eventsub

import (
	"time"
)

const (
	WelcomeMessageType      = "session_welcome"
	KeepAliveMessageType    = "session_keepalive"
	NotificationMessageType = "notification"
)

type TwitchMessage struct {
	Metadata Metadata       `json:"metadata"`
	Payload  map[string]any `json:"payload"`
}

type Metadata struct {
	MessageID        string    `json:"message_id"`
	MessageType      string    `json:"message_type"`
	MessageTimestamp time.Time `json:"message_timestamp"`
}

type Session struct {
	ID                      string    `json:"id"`
	Status                  string    `json:"status"`
	ConnectedAt             time.Time `json:"connected_at"`
	KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
	ReconnectURL            *string   `json:"reconnect_url"`
	RecoveryURL             *string   `json:"recovery_url"`
}
