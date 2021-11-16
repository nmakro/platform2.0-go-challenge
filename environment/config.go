package environment

import (
	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:31000")
	viper.SetDefault("SESSION_AUTH_KEY", "!@s+za%@DaazSq@3")
}

type Config struct {
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	SessionAuthKey string `mapstructure:"SESSION_AUTH_KEY"`
}

func LoadConfig(path string) {
	viper.AddConfigPath(path)
	setDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()

	viper.AutomaticEnv()
}
