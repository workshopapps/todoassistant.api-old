package utils

import "github.com/spf13/viper"

// Config file stores configuration of application

// Values read by viper viper from a config file or enviroment variable

type Config struct {
	DataSourceName string `mapstructure:"DATA_SOURCE_NAME"`
	SeverAddress   string `mapstructure:"SEVER_ADDRESS"`
	FromEmail      string `mapstructure:"FROM_EMAIL"`
	Password       string `mapstructure:"Password"`
	Host           string `mapstructure:"Host"`
	Port           string `mapstructure:"Port"`
	TokenSecret    string `mapstructure:"TOKEN_SECRET"`
	FromEmailAddr  string `mapstructure:"FromEmailAddr"`
	SMTPpwd		   string `mapstructure:"SMTPpwd"`
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
