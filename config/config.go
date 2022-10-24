package config

import "github.com/spf13/viper"

type Configuration struct {
	Environment string
	Token       string
	Mongo       MongoConfiguration
}

type MongoConfiguration struct {
	Server     string
	Database   string
	Collection string
}

func GetConfig() Configuration {
	var config Configuration

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}
