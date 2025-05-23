package subscriptions

import (
	"twithoauth/storage"
	"twithoauth/types"
)

type SubscriptionHandler struct {
	store *storage.Storage
}

func NewSubscriptionHandler(store *storage.Storage) *SubscriptionHandler {
	return &SubscriptionHandler{
		store: store,
	}
}

func (h *SubscriptionHandler) SubscribeChannelFollow(broadcasterID string, token types.UserAccessToken, clientID string) {

}
