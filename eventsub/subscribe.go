package eventsub

import (
	"bufio"
	"encoding/json"
	"net/http"
	"subalertor/logger"
	"subalertor/utils"
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
	tokens utils.Tokens,
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

	logger.Log.WithField("json_body", string(data)).Debug("Prepared JSON")

	res, err := utils.DFetcher.FetchTwitchApi(
		"https://api.twitch.tv/helix/eventsub/subscriptions",
		http.MethodPost,
		data,
		&utils.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	defer res.Body.Close()

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
