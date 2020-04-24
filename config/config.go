package config

import "time"

type Config struct {
	AppPort       string
	FetchInterval time.Duration
	RedisUri      string
	RedisTimeout  time.Duration
	TelegramToken string
}

func NewConfig() *Config {
	return &Config{
		AppPort:       GetString("APP_PORT"),
		FetchInterval: time.Duration(GetInt("FETCH_SCHEDULED_IN_S")) * time.Second,
		RedisUri:      GetString("REDIS_URI"),
		RedisTimeout:  time.Duration(GetInt("REDIS_TIMEOUT_IN_S")) * time.Second,
		TelegramToken: GetString("TELEGRAM_TOKEN"),
	}
}
