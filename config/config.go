package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/melody-mood/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	app struct {
		Env      string
		Name     string
		LogLevel string
	}

	http struct {
		Port int
	}

	spotify struct {
		ClientID     string
		ClientSecret string
		RedirectURI  string
	}

	openAI struct {
		APIKey string
	}

	redis struct {
		Host string
		Port string
	}

	Config struct {
		App     app
		Http    http
		Spotify spotify
		OpenAI  openAI
		Redis   redis
	}
)

var (
	configData *Config
)

func InitConfig() {
	viper.SetConfigType("env")
	viper.SetConfigName(".env") // name of Config file (without extension)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if os.Getenv("APP_ENV") == constants.DEV {
		if err := viper.ReadInConfig(); err != nil {
			logrus.WithError(err).Warn("failed to load config file")
		}
	}

	configData = &Config{
		App: app{
			Env:  getRequiredString("APP_ENV"),
			Name: getRequiredString("APP_NAME"),
		},
		Http: http{
			Port: getRequiredInt("APP_PORT"),
		},
		Spotify: spotify{
			ClientID:     getRequiredString("SPOTIFY_CLIENT_ID"),
			ClientSecret: getRequiredString("SPOTIFY_CLIENT_SECRET"),
			RedirectURI:  getRequiredString("REDIRECT_URI"),
		},
		OpenAI: openAI{
			APIKey: getRequiredString("OPENAI_API_KEY"),
		},
		Redis: redis{
			Host: getRequiredString("REDIS_HOST"),
			Port: getRequiredString("REDIS_PORT"),
		},
	}
}

func getRequiredString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	log.Fatalln(fmt.Errorf("KEY %s IS MISSING", key))
	return ""
}

func getRequiredInt(key string) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

func getRequiredBool(key string) bool {
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

func getRequiredDuration(key string) time.Duration {
	if viper.IsSet(key) {
		return viper.GetDuration(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

func GetConfig() Config {
	return *configData
}
