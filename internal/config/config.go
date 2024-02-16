package config

import (
	"ewallet/pkg/model"

	"github.com/spf13/viper"
)

func InitConfig() error {

	viper.SetConfigFile("./internal/config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return model.ErrLoadConfig.Error()
	}

	return nil
}
