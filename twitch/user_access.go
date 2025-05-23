package twitch

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"

	"twithoauth/eventsub"
	"twithoauth/logger"
	"twithoauth/types"
)

const TwitchAuthURL = "https://id.twitch.tv/oauth2/authorize"

var states = make(map[string]string)

func AuthorizeHandler(clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(TwitchAuthURL)
		if err != nil {
			logger.Log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		state := GenerateRandomState()

		query := url.Query()
		query.Set("client_id", clientID)
		query.Set("redirect_uri", "http://localhost:8080/callback")
		query.Set("response_type", "code")
		query.Set("scope", "user:read:email moderator:read:followers")
		query.Set("state", state)
		url.RawQuery = query.Encode()

		http.Redirect(w, r, url.String(), http.StatusSeeOther)
	}
}

type UserAccessTokenResponse struct {
	AccessToken  types.UserAccessToken `json:"access_token"`
	ExpiresIn    int64                 `json:"expires_in"`
	RefreshToken string                `json:"refresh_token"`
	Scope        []string              `json:"scope"`
	TokenType    string                `json:"token_type"`
}

func CallbackHandler(clientID string, clientSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		if state == "" || states[state] != state {
			logger.Log.Error("Invalid state")
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		form := url.Values{
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"code":          {code},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {"http://localhost:8080/callback"},
		}

		req, err := http.PostForm("https://id.twitch.tv/oauth2/token", form)
		if err != nil {
			logger.Log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		jsonDecoder := json.NewDecoder(req.Body)
		var token UserAccessTokenResponse
		err = jsonDecoder.Decode(&token)
		if err != nil {
			logger.Log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Log.
			WithField("expiresIn", token.ExpiresIn).
			Debug("Got access token")

		userData, err := GetUserData(token.AccessToken, clientID)
		if err != nil {
			logger.Log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		userID := userData.Data[0].ID
		logger.Log.
			WithField("user", userID).
			Debug("Got user data")

		eventsub.SubscribeChannelFollow(
			userID,
			token.AccessToken,
			clientID,
		)

		w.WriteHeader(http.StatusOK)
	}
}

func GenerateRandomState() string {
	byt := make([]byte, 32)
	rand.Read(byt)
	str := hex.EncodeToString(byt)
	states[str] = str
	return str
}
