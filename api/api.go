package api

import (
	"net/http"
	"twithoauth/logger"
	"twithoauth/storage"
	"twithoauth/twitch"
)

func Run(twitchClientGrant *twitch.ClientCredentials, storage *storage.Storage) {
	mux := http.NewServeMux()

	twitchHandler := twitch.NewTwitchHandlers(twitchClientGrant.ClientID, twitchClientGrant.ClientSecret, storage)

	mux.HandleFunc("/login", twitchHandler.AuthorizeHandler)
	mux.HandleFunc("/callback", twitchHandler.CallbackHandler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Log.Fatal(err)
	}
}
