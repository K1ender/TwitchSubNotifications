package api

import (
	"net/http"
	"twithoauth/handlers"
	"twithoauth/logger"
	"twithoauth/middleware"
	"twithoauth/storage"
	"twithoauth/twitch"
)

func Run(twitchClientGrant *twitch.ClientCredentials, storage *storage.Storage) {
	mux := http.NewServeMux()

	twitchHandler := twitch.NewTwitchHandlers(twitchClientGrant.ClientID, twitchClientGrant.ClientSecret, storage)

	profileHandler := handlers.NewProfileHandler(storage)
	authMiddleware := middleware.AuthMiddleware(storage)

	mux.HandleFunc("GET /login", twitchHandler.AuthorizeHandler)
	mux.HandleFunc("GET /callback", twitchHandler.CallbackHandler)

	mux.Handle("GET /profile", authMiddleware(http.HandlerFunc(profileHandler.GetProfile)))
	mux.Handle("GET /logout", authMiddleware(http.HandlerFunc(profileHandler.LogoutHandler)))

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Log.Fatal(err)
	}
}
