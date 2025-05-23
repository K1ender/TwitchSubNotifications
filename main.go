package main

import (
	"twithoauth/api"
	"twithoauth/config"
	"twithoauth/database"
	"twithoauth/eventsub"
	"twithoauth/storage"
	"twithoauth/twitch"
)

func main() {
	cfg := config.MustInit()

	db := database.MustInit(&cfg)

	twitchClientGrant := twitch.NewClientCredentials(
		cfg.Twitch.ClientID,
		cfg.Twitch.ClientSecret,
	)
	twitchClientGrant.GetAccessToken()
	go twitchClientGrant.UpdateAccessToken()

	storage := storage.NewStorage(db)

	go eventsub.EventSubHandler()
	go api.Run(twitchClientGrant, storage)

	select {}

}
