package twitch

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"subalertor/config"
	"subalertor/logger"
	"subalertor/storage"
	"subalertor/types"
	"subalertor/utils"
)

var states = make(map[string]string)

type Handlers struct {
	clientID     string
	clientSecret string
	storage      *storage.Storage
	cfg          *config.Config
}

func NewTwitchHandlers(clientID string, clientSecret string, storage *storage.Storage, cfg *config.Config) *Handlers {
	return &Handlers{
		clientID:     clientID,
		clientSecret: clientSecret,
		storage:      storage,
		cfg:          cfg,
	}
}

func (h *Handlers) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(utils.AuthorizeURL)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	state := GenerateRandomState()
	states[state] = state

	query := url.Query()
	query.Set("client_id", h.clientID)
	query.Set("redirect_uri", h.cfg.FrontEndURL+"/callback")
	query.Set("response_type", "code")
	query.Set("scope", config.Scopes)
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

func (h *Handlers) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if state == "" || states[state] != state {
		logger.Log.Error("Invalid state")
		utils.BadRequest(w, "Invalid state")
		return
	}

	form := url.Values{
		"client_id":     {h.clientID},
		"client_secret": {h.clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:5173/callback"},
	}

	res, err := http.PostForm(utils.AccessTokenURL, form)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Log.Error("Failed to get access token")
		utils.InternalServerError(w)
		return
	}

	jsonDecoder := json.NewDecoder(res.Body)
	var token UserAccessTokenResponse
	err = jsonDecoder.Decode(&token)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	logger.Log.
		WithField("expiresIn", token.ExpiresIn).
		Debug("Got access token")

	userData, err := GetUserData(utils.Tokens{
		AccessToken:  string(token.AccessToken),
		RefreshToken: token.RefreshToken,
	}, h.clientID)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
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
				utils.InternalServerError(w)
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
				utils.InternalServerError(w)
				return
			}
		} else {
			logger.Log.Error(err)
			utils.InternalServerError(w)
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
			utils.InternalServerError(w)
			return
		}
	}

	userAuthToken := GenerateRandomState()
	utils.SetAuthCookie(w, r, userAuthToken)

	if err := h.storage.SessionStore.CreateSession(
		user.ID,
		utils.HashToken(userAuthToken),
		time.Now().Add(time.Hour*24*30),
	); err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	delete(states, state)

	utils.OK(w, user.Username)
}

func GenerateRandomState() string {
	byt := make([]byte, 32)
	rand.Read(byt)
	str := hex.EncodeToString(byt)
	return str
}
