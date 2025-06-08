package websocket

import (
	"net/http"
	"subalertor/logger"
	"subalertor/storage"
	"subalertor/utils"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketMap struct {
	WebSockets map[string]*websocket.Conn
	sync.Mutex
}

var WebSockets = WebSocketMap{
	WebSockets: make(map[string]*websocket.Conn),
}

func GetUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}

type EventSubHandler struct {
	store *storage.Storage
}

func NewEventSubHandler(store *storage.Storage) *EventSubHandler {
	return &EventSubHandler{
		store: store,
	}
}

type Event[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type NewSubscriberEvent struct {
	Username string `json:"username"`
}

func (h *EventSubHandler) FollowHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		logger.Log.Error("Missing id")
		utils.BadRequest(w, "Missing id")
		return
	}

	user, err := h.store.UserStore.FindUserByID(userID)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	events, err := h.store.EventSubStore.GetSubscribedEvents(user.ID)
	if err != nil {
		logger.Log.Error(err)
		utils.InternalServerError(w)
		return
	}

	founded := false
	for _, event := range events {
		if event.Type == "channel.follow" {
			founded = true
			break
		}
	}

	if !founded {
		logger.Log.Error("User is not subscribed to channel follow")
		utils.BadRequest(w, "User is not subscribed to channel follow")
		return
	}

	ws, err := GetUpgrader().Upgrade(w, r, nil)
	if err != nil {
		return
	}
	logger.Log.Debug("Upgraded connection")
	WebSockets.Lock()
	WebSockets.WebSockets[userID] = ws
	WebSockets.Unlock()
}
