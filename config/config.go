package config

import (
	"bytes"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var defaultConfig []byte

type Database struct {
	Host     string
	Port     uint
	Username string
	Password string
	Dbname   string
}

type Server struct {
	Port string
}

type Config struct {
	Database Database
	Server   Server
	Policy   Policy
}

type Policy struct {
	Users KafkaConfig
}

type KafkaConfig struct {
	Topic            string
	BootStrapServers string
	GroupId          string
	Timeout          int
}

func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		panic(err)
	}
	return config
}

func Load() (*Config, error) {
	// viper
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// path
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
		return nil, err
	}

	host := viper.GetString("database.host")
	log.Printf("host from config: %s", host)

	// if host is empty, use from config.yaml
	if host == "" {
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
