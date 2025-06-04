package api

import (
	"net/http"
	"subalertor/config"
	"subalertor/handlers"
	"subalertor/logger"
	"subalertor/middleware"
	"subalertor/storage"
	"subalertor/twitch"
)

func Run(twitchClientGrant *twitch.ClientCredentials, storage *storage.Storage, cfg *config.Config) {
	mux := http.NewServeMux()

	twitchHandler := twitch.NewTwitchHandlers(twitchClientGrant.ClientID, twitchClientGrant.ClientSecret, storage, cfg)

	profileHandler := handlers.NewProfileHandler(storage)
	subscriptionHandler := handlers.NewSubscriptionHandler(storage, cfg)
	authMiddleware := middleware.AuthMiddleware(storage)
	corsMiddleware := middleware.CORS(cfg)

	mux.HandleFunc("GET /login", twitchHandler.AuthorizeHandler)
	mux.HandleFunc("GET /callback", twitchHandler.CallbackHandler)

	mux.Handle("GET /profile", authMiddleware(http.HandlerFunc(profileHandler.GetProfile)))
	mux.Handle("POST /logout", authMiddleware(http.HandlerFunc(profileHandler.LogoutHandler)))

	mux.Handle("GET /followers", authMiddleware(http.HandlerFunc(profileHandler.GetLatestFollowers)))

	mux.Handle("POST /subscribe/{channel_id}", authMiddleware(http.HandlerFunc(subscriptionHandler.SubscribeChannelFollowHandler)))
	mux.Handle("POST /unsubscribe/{id}", authMiddleware(http.HandlerFunc(subscriptionHandler.UnsubscribeChannelFollowHandler)))

	mux.Handle("GET /subscribed", authMiddleware(http.HandlerFunc(profileHandler.GetSubscribedEvents)))

	wrapped := middleware.Use(mux, corsMiddleware)

	if err := http.ListenAndServe(":8080", wrapped); err != nil {
		logger.Log.Fatal(err)
	}
}
