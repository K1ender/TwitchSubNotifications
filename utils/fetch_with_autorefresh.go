package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"twithoauth/logger"
	"twithoauth/storage"
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

func NewFetcher(clientID string, tokenStorage storage.TokenStore) *Fetcher {
	return &Fetcher{
		ClientID:     clientID,
		TokenStorage: tokenStorage,
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
	}

	req, err := http.PostForm(AccessTokenURL, form)
	if err != nil {
		logger.Log.Error(err)
		return "", "", err
	}

	jsonDecoder := json.NewDecoder(req.Body)
	var token UpdatedToken
	err = jsonDecoder.Decode(&token)
	if err != nil {
		logger.Log.Error(err)
		return "", "", err
	}
	return token.AccessToken, token.RefreshToken, nil
}

func (f *Fetcher) FetchTwitchApi(req *http.Request, tokens *Tokens) (*http.Response, error) {
	req.Header.Add("Client-Id", f.ClientID)
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	if res.StatusCode == 401 {
		accessToken, refreshToken, err := f.RefreshAccessUserToken(tokens.RefreshToken)
		if err != nil {
			logger.Log.Error(err)
			return nil, err
		}
		tokens.AccessToken = accessToken
		tokens.RefreshToken = refreshToken
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		return f.FetchTwitchApi(req, tokens)
	}

	return res, nil
}
