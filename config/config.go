package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BotToken         string
	OpenWeatherToken string
}

var cfg *Config

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	cfg = &Config{}
}

func Read() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	cfg.BotToken = viper.GetString("BOT_SECRET_KEY")
	cfg.OpenWeatherToken = viper.GetString("OPEN_WEATHER_KEY")

	return nil
}

func GetConfig() *Config {
	return cfg
}
