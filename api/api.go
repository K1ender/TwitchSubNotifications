package api

import (
	"net/http"
	"sync"
	"twithoauth/twitch"
)

func Run(wg *sync.WaitGroup, twitchClientGrant *twitch.ClientCredentials) {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", twitch.AuthorizeHandler(twitchClientGrant.ClientID))
	mux.HandleFunc("/callback", twitch.CallbackHandler(twitchClientGrant.ClientID, twitchClientGrant.ClientSecret))

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	wg.Done()
	srv.ListenAndServe()
}
