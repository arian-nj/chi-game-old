package config

import (
	"fmt"
	"os"
)

type ReleasMode int

const (
	Release ReleasMode = 0
	Develop ReleasMode = 1
)

type JWTConfig struct {
	SecretKey []byte
}

type Config struct {
	BotToken    string
	ReleaseMode ReleasMode
	DatabseUrl  string
	BaseUrl     string
	Jwt         JWTConfig
}

func ParseConfig() (*Config, error) {
	var conf Config

	release_mode := os.Getenv("RELEASE_MODE")
	if release_mode == "" {
		return nil, fmt.Errorf("env var release mode is empty")
	}

	if release_mode == "release" {
		conf.ReleaseMode = Release

		bot_token := os.Getenv("RELEASE_BOT_TOKEN")
		if bot_token == "" {
			return nil, fmt.Errorf("release bot token is empty")
		}

		db_url := os.Getenv("RELEASE_DATABASE_URL")
		if db_url == "" {
			return nil, fmt.Errorf("release database url is empy")
		}
		conf.BotToken = bot_token
		conf.DatabseUrl = db_url

	} else {
		conf.ReleaseMode = Develop

		bot_token := os.Getenv("DEV_BOT_TOKEN")
		if bot_token == "" {
			return nil, fmt.Errorf("release bot token is empty")
		}

		db_url := os.Getenv("DEV_DATABASE_URL")
		if db_url == "" {
			return nil, fmt.Errorf("release database url is empy")
		}

		conf.BotToken = bot_token
		conf.DatabseUrl = db_url
	}
	conf.Jwt = JWTConfig{
		SecretKey: []byte("palfd34nm3n892n74riq3v4k235vq65b5q7n7nnw7nwsp8w5b4bq34v5bqb56q35v5b6n7opawrr"),
	}
	conf.BaseUrl = "https://192.168.0.107"
	return &conf, nil
}
