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
	authMiddleware := middleware.AuthMiddleware(storage)
	corsMiddleware := middleware.CORS(cfg)

	mux.HandleFunc("GET /login", twitchHandler.AuthorizeHandler)
	mux.HandleFunc("GET /callback", twitchHandler.CallbackHandler)

	mux.Handle("GET /profile", authMiddleware(http.HandlerFunc(profileHandler.GetProfile)))
	mux.Handle("GET /logout", authMiddleware(http.HandlerFunc(profileHandler.LogoutHandler)))

	mux.Handle("GET /followers", authMiddleware(http.HandlerFunc(profileHandler.GetLatestFollowers)))

	wrapped := middleware.Use(mux, corsMiddleware)

	if err := http.ListenAndServe(":8080", wrapped); err != nil {
		logger.Log.Fatal(err)
	}
}
