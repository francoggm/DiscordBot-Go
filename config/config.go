package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BotToken    string
	AppID       string
	GuildID     string
	StockAPIKey string
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
	cfg.AppID = viper.GetString("APP_ID")
	cfg.GuildID = viper.GetString("GUILD_ID")
	cfg.StockAPIKey = viper.GetString("STOCK_API")

	return nil
}

func GetConfig() *Config {
	return cfg
}
