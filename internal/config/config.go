package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

var cfg Config

type Config struct {
	Postgres              DatabaseConfig `yaml:"DB"`
	HistoryPostgres       DatabaseConfig `yaml:"History_DB"`
	RedisProfilesToUpdate RedisConfig    `yaml:"Redis_profiles_to_update"`
	RedisUpdatedProfiles  RedisConfig    `yaml:"Redis_updated_profiles"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type RedisConfig struct {
	Broker          string        `yaml:"broker"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Password        string        `yaml:"password"`
	PollingInterval time.Duration `yaml:"polling_interval"`
	DeliveryMethod  string        `yaml:"deliveryMethod"`
	Topic           string        `yaml:"topic"`
}

func InitConfig() (*Config, error) {
	configPath := "cmd/config.yaml"

	clean := filepath.Clean(configPath)

	file, err := os.Open(clean)
	if err != nil {
		return nil, fmt.Errorf("fail to open config file in path \"%s\" with error %w", configPath, err)
	}

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config %w", err)
	}

	return &cfg, nil
}

func GetConfig() *Config {
	return &cfg
}
