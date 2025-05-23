package twitch

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"twithoauth/logger"
	"twithoauth/storage"
	"twithoauth/types"
	"twithoauth/utils"
)

const TwitchAuthURL = "https://id.twitch.tv/oauth2/authorize"

var states = make(map[string]string)

type TwitchHandlers struct {
	clientID     string
	clientSecret string
	storage      *storage.Storage
}

const scopes = "user:read:email moderator:read:followers"

func NewTwitchHandlers(clientID string, clientSecret string, storage *storage.Storage) *TwitchHandlers {
	return &TwitchHandlers{
		clientID:     clientID,
		clientSecret: clientSecret,
		storage:      storage,
	}
}

func (h *TwitchHandlers) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(TwitchAuthURL)
	if err != nil {
		logger.Log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	state := GenerateRandomState()
	states[state] = state

	query := url.Query()
	query.Set("client_id", h.clientID)
	query.Set("redirect_uri", "http://localhost:8080/callback")
	query.Set("response_type", "code")
	query.Set("scope", scopes)
	query.Set("state", state)
	url.RawQuery = query.Encode()

	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}

type UserAccessTokenResponse struct {
	AccessToken  types.UserAccessToken `json:"access_token"`
	ExpiresIn    int64                 `json:"expires_in"`
	RefreshToken string                `json:"refresh_token"`
	Scope        []string              `json:"scope"`
	TokenType    string                `json:"token_type"`
}

func (h *TwitchHandlers) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if state == "" || states[state] != state {
		logger.Log.Error("Invalid state")
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	form := url.Values{
		"client_id":     {h.clientID},
		"client_secret": {h.clientSecret},
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

	userData, err := GetUserData(token.AccessToken, h.clientID)
	if err != nil {
		logger.Log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userID := userData.Data[0].ID
	logger.Log.
		WithField("user", userID).
		Debug("Got user data")

	user, err := h.storage.UserStore.CreateUser(userID, userData.Data[0].Login)
	if err != nil {
		if err == storage.ErrUserAlreadyExists {
			user, err = h.storage.UserStore.FindUserByID(userID)
			if err != nil {
				logger.Log.Error(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			err = h.storage.TokenStore.SetTokens(
				user.ID,
				string(token.AccessToken),
				token.RefreshToken,
				time.Unix(token.ExpiresIn, 0),
			)
			if err != nil {
				logger.Log.Error(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			logger.Log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.storage.TokenStore.AddTokens(
			user.ID,
			string(token.AccessToken),
			token.RefreshToken,
			time.Unix(token.ExpiresIn, 0),
		); err != nil {
			logger.Log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	userAuthToken := GenerateRandomState()
	utils.SetAuthCookie(w, r, userAuthToken)

	if err := h.storage.SessionStore.CreateSession(user.ID, userAuthToken, time.Now().Add(time.Hour*24*30)); err != nil {
		logger.Log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	delete(states, state)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func GenerateRandomState() string {
	byt := make([]byte, 32)
	rand.Read(byt)
	str := hex.EncodeToString(byt)
	return str
}

type UpdatedToken struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func RefreshAccessUserToken(clientID string, clientSecret string, refreshToken string) (accessToken string, refresh_token string, err error) {
	form := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, err := http.PostForm("https://id.twitch.tv/oauth2/token", form)
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
