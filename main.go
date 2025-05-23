package main

import (
	"sync"
	"twithoauth/api"
	"twithoauth/config"
	"twithoauth/eventsub"
	"twithoauth/twitch"
)

func main() {
	cfg := config.MustInit()

	twitchClientGrant := twitch.NewClientCredentials(
		cfg.Twitch.ClientID,
		cfg.Twitch.ClientSecret,
	)
	twitchClientGrant.GetAccessToken()
	go twitchClientGrant.UpdateAccessToken()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Add(1)
	go eventsub.EventSubHandler(&wg)
	go api.Run(&wg, twitchClientGrant)

	wg.Wait()
	select {}

}
