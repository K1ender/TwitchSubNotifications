package twitch

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"twithoauth/logger"
	"twithoauth/utils"
)

type ClientCredentials struct {
	ClientID     string
	ClientSecret string
	AccessToken  AppAccessToken
	ExpiresIn    int
}

func NewClientCredentials(clientID string, clientSecret string) *ClientCredentials {
	return &ClientCredentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AppAccessToken string

func (c *ClientCredentials) GetAccessToken() (AppAccessToken, error) {
	form := url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"client_credentials"},
	}
	resp, err := http.PostForm(utils.AccessTokenURL, form)
	if err != nil {
		logger.Log.Error(err)
		return "", err
	}

	decoder := json.NewDecoder(resp.Body)
	var token Token
	err = decoder.Decode(&token)
	if err != nil {
		logger.Log.Error(err)
		return "", err
	}
	logger.Log.
		WithField("expiresIn", token.ExpiresIn).
		Debug("Got access token")
	c.AccessToken = AppAccessToken(token.AccessToken)
	c.ExpiresIn = token.ExpiresIn

	return c.AccessToken, nil
}

func (c *ClientCredentials) GetBearerToken() string {
	return "Bearer " + string(c.AccessToken)
}

func (c *ClientCredentials) UpdateAccessToken() {
	for {
		time.Sleep(time.Duration(c.ExpiresIn-60) * time.Second)
		c.GetAccessToken()
	}
}
