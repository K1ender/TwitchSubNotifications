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

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(user.Username))
}

func (h *ProfileHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetAuthCookie(r)
	if err != nil {
		logger.Log.Error(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	h.store.SessionStore.DeleteSession(utils.HashToken(token))
	utils.DeleteAuthCookie(w, r)

	w.WriteHeader(http.StatusOK)
}
