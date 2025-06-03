package eventsub

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"subalertor/logger"
	"subalertor/utils"
)

const EventApiEndpoint = "https://api.twitch.tv/helix/eventsub/subscriptions"

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

type Response struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func SubscribeChannelFollow(
	broadcasterID string,
	tokens utils.Tokens,
	clientID string,
) (Response, error) {
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
		return Response{}, err
	}

	logger.Log.WithField("json_body", string(data)).Debug("Prepared JSON")

	res, err := utils.DFetcher.FetchTwitchApi(
		EventApiEndpoint,
		http.MethodPost,
		data,
		&utils.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	)
	if err != nil {
		logger.Log.Error(err)
		return Response{}, err
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.Log.Error(err)
		}
	}(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		response, _ := bufio.NewReader(res.Body).ReadString('\n')
		logger.Log.
			WithField("statusCode", res.StatusCode).
			WithField("errorResponse", response).
			Error("Failed to subscribe")
		return Response{}, fmt.Errorf("failed to subscribe")
	}

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		logger.Log.Error(err)
		return Response{}, err
	}
	return response, nil
}

func UnsubscribeChannelFollow(subscriptionID, clientID string, tokens utils.Tokens) error {
	url, err := url.Parse(EventApiEndpoint)
	if err != nil {
		logger.Log.Error(err)
		return err
	}
	query := url.Query()
	query.Add("id", subscriptionID)
	url.RawQuery = query.Encode()

	res, err := utils.DFetcher.FetchTwitchApi(
		url.String(),
		http.MethodDelete,
		nil,
		&utils.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.Log.Error(err)
		}
	}(res.Body)

	if res.StatusCode != 200 && res.StatusCode != 204 {
		logger.Log.Error("Failed to unsubscribe")
		return err
	}
	logger.Log.Info("Unsubscribed from channel follow", subscriptionID)
	return nil
}
