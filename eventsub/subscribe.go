package eventsub

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"twithoauth/logger"
	"twithoauth/types"
)

type CreateEventSub[T any] struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition T         `json:"condition"`
	Transport Transport `json:"transport"`
}

type Transport struct {
	Method    string        `json:"method"`
	SessionID SessionIDType `json:"session_id"`
}

type ChannelFollowSubscription struct {
	BroadCasterUserID string `json:"broadcaster_user_id"`
	ModeratorUserID   string `json:"moderator_user_id"`
}

func SubscribeChannelFollow(
	broadcasterID string,
	token types.UserAccessToken,
	clientID string,
) error {
	body := CreateEventSub[ChannelFollowSubscription]{
		Type:    "channel.follow",
		Version: "2",
		Condition: ChannelFollowSubscription{
			BroadCasterUserID: broadcasterID,
			ModeratorUserID:   broadcasterID,
		},
		Transport: Transport{
			Method:    "websocket",
			SessionID: SessionID,
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewReader(data))
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+string(token))
	req.Header.Add("Client-Id", clientID)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 && res.StatusCode != 202 {
		response, _ := bufio.NewReader(res.Body).ReadString('\n')
		logger.Log.
			WithField("statusCode", res.StatusCode).
			WithField("errorResponse", response).
			Error("Failed to subscribe")
		return err
	}
	return nil
}
