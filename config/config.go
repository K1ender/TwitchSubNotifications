package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Twitch TwitchConfig
}

type TwitchConfig struct {
	ClientID     string `env:"TWITCH_CLIENT_ID" env-required:"true"`
	ClientSecret string `env:"TWITCH_CLIENT_SECRET" env-required:"true"`
}

func MustInit() Config {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
