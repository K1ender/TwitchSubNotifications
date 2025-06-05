package main

import (
	"subalertor/api"
	"subalertor/config"
	"subalertor/database"
	"subalertor/eventsub"
	"subalertor/logger"
	"subalertor/storage"
	"subalertor/twitch"
	"subalertor/utils"
	"sync"
)

func main() {
	cfg := config.MustInit()

	db := database.MustInit(&cfg)

	twitchClientGrant := twitch.NewClientCredentials(
		cfg.Twitch.ClientID,
		cfg.Twitch.ClientSecret,
	)
	_, err := twitchClientGrant.GetAccessToken()
	if err != nil {
		logger.Log.Fatal(err)
	}
	go twitchClientGrant.UpdateAccessToken()

	storage := storage.NewStorage(db)

	utils.DFetcher = utils.NewFetcher(
		cfg.Twitch.ClientID,
		cfg.Twitch.ClientSecret,
		storage.TokenStore,
	)

	var wg sync.WaitGroup
	wg.Add(1)

	go eventsub.EventSubHandler(&wg, storage)
	go api.Run(twitchClientGrant, storage, &cfg)

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

			accessToken, refreshToken, err := storage.TokenStore.GetTokens(*userID)
			if err != nil {
				continue
			}

			_, err = eventsub.SubscribeChannelFollow(
				*event.Condition.BroadcasterID,
				utils.Tokens{
					AccessToken:  string(accessToken),
					RefreshToken: refreshToken,
				},
				twitchClientGrant.ClientID,
			)
			if err != nil {
				logger.Log.Error(err)
				continue
			}
			logger.Log.Info("Subscribed to channel follow")
		}
	}

	select {}
}
