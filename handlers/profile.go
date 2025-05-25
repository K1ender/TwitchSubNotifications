package handlers

import (
	"net/http"
	"subalertor/logger"
	"subalertor/middleware"
	"subalertor/storage"
	"subalertor/utils"
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
	err = h.store.SessionStore.DeleteSession(utils.HashToken(token))
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}
	utils.DeleteAuthCookie(w, r)

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Success",
	})
}
