package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
}

func (config AppConfig) LoadEnv() {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Error loading .env file:", err)

	}

	keys := viper.AllKeys()

	log.Println(keys)

	for _, key := range keys {
		key = strings.ToUpper(key)
		log.Println(key)
		log.Println(viper.GetString(key))

		envValue, found := os.LookupEnv(key)
		if !found || envValue == "" {
			os.Setenv(key, viper.GetString(key))
		}

	}
}

func (config AppConfig) SetLogLevel() {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		log.Warn("LOG_LEVEL is not set, defaulting to 'info'")
		levelStr = "info"
	}

}
