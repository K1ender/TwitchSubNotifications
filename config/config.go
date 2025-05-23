package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Twitch      TwitchConfig
	Database    DatabaseConfig
	FrontEndURL string `env:"FRONTEND_URL" env-required:"true"`
}

type TwitchConfig struct {
	ClientID     string `env:"TWITCH_CLIENT_ID" env-required:"true"`
	ClientSecret string `env:"TWITCH_CLIENT_SECRET" env-required:"true"`
}

type DatabaseConfig struct {
	File string `env:"DATABASE_FILE" env-default:"./database.db"`
}

func MustInit() Config {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
