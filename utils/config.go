package utils

import "github.com/spf13/viper"

type Config struct {

	DbDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	DBName string `mapstructure:"DB_NAME"`
	SigningKey string `mapstructure:"SIGNING_KEY"`
}


func LoadConfig(path string) (config *Config, erorr error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("env")
	viper.SetConfigType("env")

	viper.AutomaticEnv();

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}