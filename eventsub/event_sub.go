package eventsub

import (
	"subalertor/logger"
	"subalertor/storage"
	events "subalertor/websocket"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var SessionID SessionIDType

type SessionIDType string

type Event[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type NewSubscriberEvent struct {
	Username string `json:"username"`
}

func EventSubHandler(wg *sync.WaitGroup, store *storage.Storage) {
	ws, _, err := websocket.DefaultDialer.Dial("wss://eventsub.wss.twitch.tv/ws", nil)
	if err != nil {
		panic(err)
	}

	for {
		var twitchMessage TwitchMessage
		err := ws.ReadJSON(&twitchMessage)
		if err != nil {
			logger.Log.Error(err)
			break
		}

		if twitchMessage.Metadata.MessageType == WelcomeMessageType {
			data, ok := twitchMessage.Payload["session"].(map[string]any)
			if !ok {
				logger.Log.
					WithField("twitchMessage", twitchMessage).
					Error("Failed to get sessionID")
				continue
			}
			id, ok := data["id"].(string)
			if !ok {
				logger.Log.
					WithField("twitchMessage", twitchMessage).
					Error("Failed to get sessionID")
				continue
			}
			SessionID = SessionIDType(id)
			logger.Log.
				WithField("sessionID", SessionID).
				Info("✅ Connected to Twitch Eventsub API")
			wg.Done()
		} else if twitchMessage.Metadata.MessageType == KeepAliveMessageType {

		} else if twitchMessage.Metadata.MessageType == NotificationMessageType {
			subscription := twitchMessage.Payload["subscription"].(map[string]any)
			subscriptionType := subscription["type"].(string)
			event := twitchMessage.Payload["event"].(map[string]any)

			if subscriptionType == "channel.follow" {
				display_username := event["user_name"].(string)
				username := event["user_login"].(string)
				user_id := event["user_id"].(string)
				broadcaster_id := event["broadcaster_user_id"].(string)
				err := store.FollowerStore.AddFollower(
					broadcaster_id,
					storage.FollowerModel{
						ID:          user_id,
						DisplayName: display_username,
						Username:    username,
						FollowedAt:  int(time.Now().Unix()),
					},
				)
				if err != nil {
					logger.Log.
						WithField("error", err).
						Error("Failed to add follower")
				}
				logger.Log.
					WithField("username", username).
					Info("New follower")
				events.WebSockets.WebSockets[broadcaster_id].WriteJSON(Event[NewSubscriberEvent]{
					Type: "new_subscriber",
					Data: NewSubscriberEvent{
						Username: display_username,
					},
				})
			}
		} else {
			logger.Log.
				WithField("messageType", twitchMessage.Metadata.MessageType).
				WithField("twitchMessage", twitchMessage).
				Debug("twitchMessage")
		}
	}
}
