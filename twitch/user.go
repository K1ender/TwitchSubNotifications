package twitch

import (
	"encoding/json"
	"net/http"
	"time"
	"twithoauth/logger"
	"twithoauth/types"
)

type User struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	Type            string    `json:"type"`
	BroadcasterType string    `json:"broadcaster_type"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	OfflineImageURL string    `json:"offline_image_url"`
	ViewCount       int64     `json:"view_count"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
}

func GetUserData(accessToken types.UserAccessToken, clientID string) (types.Response[User], error) {
	var user types.Response[User]
	req, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		logger.Log.Error(err)
		return user, err
	}

	req.Header.Add("Authorization", "Bearer "+string(accessToken))
	req.Header.Add("Client-Id", clientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Error(err)
		return user, err
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		logger.Log.Error(err)
		return user, err
	}
	return user, nil
}
