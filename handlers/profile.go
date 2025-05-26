package handlers

import (
	"net/http"
	"strconv"
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

func (h *ProfileHandler) GetLatestFollowers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		logger.Log.Error(err)
		utils.BadRequest(w, "Invalid offset")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		logger.Log.Error(err)
		utils.BadRequest(w, "Invalid limit")
		return
	}
	if offset < 0 || limit < 0 {
		logger.Log.Error("Offset and limit must be greater than 0")
		utils.BadRequest(w, "Offset and limit must be greater than 0")
		return
	}

	followers, err := h.store.FollowerStore.GetFollowers(user.ID, offset, limit)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Success",
		Data:    followers,
	})
}
