package handlers

import (
	"net/http"
	"twithoauth/logger"
	"twithoauth/middleware"
	"twithoauth/storage"
	"twithoauth/utils"
)

type ProfileHandler struct {
	store *storage.Storage
}

func NewProfileHandler(store *storage.Storage) *ProfileHandler {
	return &ProfileHandler{
		store: store,
	}
}

type UserResponse struct {
	Username string `json:"username"`
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Success",
		Data: UserResponse{
			Username: user.Username,
		},
	})
}

func (h *ProfileHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetAuthCookie(r)
	if err != nil {
		logger.Log.Error(err)
		utils.Unauthorized(w)
		return
	}
	h.store.SessionStore.DeleteSession(utils.HashToken(token))
	utils.DeleteAuthCookie(w, r)

	w.WriteHeader(http.StatusOK)
}
