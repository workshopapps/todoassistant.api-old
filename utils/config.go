package utils

import "github.com/spf13/viper"

// Config file stores configuration of application

// Values read by viper viper from a config file or enviroment variable

type Config struct {
	DataSourceName string `mapstructure:"DATA_SOURCE_NAME"`
	SeverAddress   string `mapstructure:"SEVER_ADDRESS"`
	StripeKey      string `mapstructure:"Stripe_Key"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
