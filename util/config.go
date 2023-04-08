package util

import "github.com/spf13/viper"

type Config struct {
	PostgresInfo string `mapstructure:"POSTGRES_INFO"`
	Addr         string `mapstructure:"ADDR"`
	SecretKey    string `mapstructure:"SECRET_KEY"`
}

func NewConfig(path string) (*Config, error) {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(path)

	// Tell viper the name of your file
	viper.SetConfigName("app")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var c *Config
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
