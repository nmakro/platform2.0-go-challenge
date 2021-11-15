package environment

import (
	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:31000")
}

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) {
	viper.AddConfigPath(path)
	setDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()

	viper.AutomaticEnv()
}
