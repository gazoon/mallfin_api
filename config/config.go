package config

import (
	"encoding/json"
	"io/ioutil"

	"sync"

	"github.com/pkg/errors"
)

var (
	config *Config
	once   sync.Once
)

func GetConfig() *Config {
	if config == nil {
		panic("Config has not initialized yet")
	}
	return config
}

func Postgres() *PostgresSettings {
	conf := GetConfig()
	return conf.Postgres
}

func Redis() *RedisSettings {
	conf := GetConfig()
	return conf.Redis
}

func Debug() bool {
	conf := GetConfig()
	return conf.Debug
}

func Port() int {
	conf := GetConfig()
	return conf.Port
}

func AccessLog() bool {
	conf := GetConfig()
	return conf.AccessLog
}

func ServiceName() string {
	conf := GetConfig()
	return conf.ServiceName
}

func ServerID() string {
	conf := GetConfig()
	return conf.ServerID
}

func LogLevel() string {
	conf := GetConfig()
	return conf.LogLevel
}

type PostgresSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"name"`
	PoolSize int    `json:"pool_size"`
	Timeout  int    `json:"timeout"`
	Retries  int    `json:"retries"`
}
type RedisSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
type Config struct {
	LogLevel    string            `json:"log_level"`
	ServiceName string            `json:"service_name"`
	ServerID    string            `json:"server_id"`
	Debug       bool              `json:"debug"`
	Port        int               `json:"port"`
	AccessLog   bool              `json:"access_log"`
	Postgres    *PostgresSettings `json:"postgres"`
	Redis       *RedisSettings    `json:"redis"`
}

func CreateConfig(path string) *Config {
	if path == "" {
		panic("Empty config path")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Errorf("Cannot read config: %s", err))
	}
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(errors.Errorf("Cannot parse config: %s", err))
	}
	return config
}

func Initialization(configPath string) {
	once.Do(func() {
		newConfig := CreateConfig(configPath)
		config = newConfig
	})
}
