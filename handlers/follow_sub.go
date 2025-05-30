package handlers

import (
	"net/http"
	"subalertor/config"
	"subalertor/eventsub"
	"subalertor/logger"
	"subalertor/middleware"
	"subalertor/storage"
	"subalertor/utils"
)

type SubscriptionHandler struct {
	store *storage.Storage
	cfg   *config.Config
}

func NewSubscriptionHandler(store *storage.Storage, cfg *config.Config) *SubscriptionHandler {
	return &SubscriptionHandler{
		store: store,
		cfg:   cfg,
	}
}

func (h *SubscriptionHandler) SubscribeChannelFollowHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	broadcasterID := r.PathValue("channel_id")

	accessToken, refreshToken, err := h.store.TokenStore.GetTokens(user.ID)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}
	if accessToken == "" || refreshToken == "" {
		logger.Log.Error("Access token or refresh token not found")
		utils.Unauthorized(w)
		return
	}

	err = h.store.EventSubStore.AddEventSubscription(storage.EventSubModel{
		Condition: storage.ConditionModel{
			BroadcasterID: &broadcasterID,
			UserID:        &broadcasterID,
		},
		Type: "channel.follow",
	})
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	err = eventsub.SubscribeChannelFollow(broadcasterID, utils.Tokens{
		AccessToken:  string(accessToken),
		RefreshToken: refreshToken,
	}, user.ID)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Success",
	})
}
