package main

import (
	"sync"
	"twithoauth/api"
	"twithoauth/config"
	"twithoauth/database"
	"twithoauth/eventsub"
	"twithoauth/logger"
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

	var wg sync.WaitGroup
	wg.Add(1)

	go eventsub.EventSubHandler(&wg)
	go api.Run(twitchClientGrant, storage)

	wg.Wait()
	events, err := storage.EventSubStore.GetAllEventSubscriptions()
	if err != nil {
		logger.Log.Fatal(err)
	}

	for _, event := range events {
		logger.Log.Debug("Event", event.PrettyPrint())
		if event.Type == "channel.follow" {
			userID := event.Condition.UserID
			if userID == nil {
				continue
			}

			accessToken, _, err := storage.TokenStore.GetTokens(*userID)
			if err != nil {
				continue
			}

			eventsub.SubscribeChannelFollow(
				*event.Condition.BroadcasterID,
				accessToken,
				twitchClientGrant.ClientID,
			)
			logger.Log.Info("Subscribed to channel follow")
		}
	}

	select {}

}
