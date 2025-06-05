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

	res, err := eventsub.SubscribeChannelFollow(broadcasterID, utils.Tokens{
		AccessToken:  string(accessToken),
		RefreshToken: refreshToken,
	}, user.ID)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	logger.Log.WithField("subscription_id", res.Data[0].ID).Debug("Subscribed to channel follow")
	err = h.store.EventSubStore.AddEventSubscription(res.Data[0].ID, storage.EventSubModel{
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

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Success",
	})
}

func (h *SubscriptionHandler) UnsubscribeChannelFollowHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	id := r.PathValue("id")

	err := h.store.EventSubStore.DeleteEventSubscription(id)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	accessToken, refreshToken, err := h.store.TokenStore.GetTokens(user.ID)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	err = eventsub.UnsubscribeChannelFollow(id, h.cfg.Twitch.ClientID, utils.Tokens{
		AccessToken:  string(accessToken),
		RefreshToken: refreshToken,
	})
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	logger.Log.Info("Unsubscribed from channel follow", id)

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Success",
	})
}
