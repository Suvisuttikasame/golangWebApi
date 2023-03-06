package util

import "github.com/spf13/viper"

type config struct {
	PostgresInfo string `mapstructure:"POSTGRES_INFO"`
	Addr         string `mapstructure:"ADDR"`
}

func NewConfig() (*config, error) {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName("app")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var c *config
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
