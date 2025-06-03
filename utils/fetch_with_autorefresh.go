package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"subalertor/config"
	"subalertor/logger"
	"subalertor/storage"
)

const (
	AuthorizeURL   = "https://id.twitch.tv/oauth2/authorize"
	AccessTokenURL = "https://id.twitch.tv/oauth2/token"
	HelixURL       = "https://api.twitch.tv/helix"
)

var DFetcher *Fetcher

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Fetcher struct {
	ClientID     string
	ClientSecret string
	TokenStorage storage.TokenStore
}

func NewFetcher(clientID string, clientSecret string, tokenStorage storage.TokenStore) *Fetcher {
	return &Fetcher{
		ClientID:     clientID,
		TokenStorage: tokenStorage,
		ClientSecret: clientSecret,
	}
}

type UpdatedToken struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func (f *Fetcher) RefreshAccessUserToken(refreshToken string) (accessToken string, refresh_token string, err error) {
	form := url.Values{
		"client_id":     {f.ClientID},
		"client_secret": {f.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"scopes":        {config.Scopes},
	}

	res, err := http.PostForm(AccessTokenURL, form)
	if err != nil {
		logger.Log.Error(err)
		return "", "", err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.Log.Error(err)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		logger.Log.Error("Failed to refresh access token")
		return "", "", fmt.Errorf("failed to refresh access token")
	}

	jsonDecoder := json.NewDecoder(res.Body)
	var token UpdatedToken
	err = jsonDecoder.Decode(&token)
	if err != nil {
		logger.Log.Error(err)
		return "", "", err
	}

	err = f.TokenStorage.UpdateAccessToken(token.RefreshToken, token.AccessToken)
	if err != nil {
		logger.Log.Error(err)
		return "", "", err
	}

	return token.AccessToken, token.RefreshToken, nil
}

func (f *Fetcher) FetchTwitchApi(url string, method string, body []byte, tokens *Tokens) (*http.Response, error) {
	makeRequest := func(token string) (*http.Response, error) {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Client-Id", f.ClientID)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		return http.DefaultClient.Do(req)
	}

	res, err := makeRequest(tokens.AccessToken)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 400 {
		logger.Log.Error("Bad request")
		return nil, fmt.Errorf("bad request")
	}

	if res.StatusCode == 401 {
		logger.Log.Debug("Access token expired. Refreshing...")
		accessToken, refreshToken, err := f.RefreshAccessUserToken(tokens.RefreshToken)
		if err != nil {
			logger.Log.Error(err)
			return nil, err
		}
		tokens.AccessToken = accessToken
		tokens.RefreshToken = refreshToken
		return makeRequest(tokens.AccessToken)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		logger.Log.Error("Failed to fetch Twitch API", res.StatusCode)
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			logger.Log.Error(scanner.Text())
		}
		return nil, fmt.Errorf("failed to fetch Twitch API")
	}

	return res, nil
}
