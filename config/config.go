package config

import (
	"encoding/json"
	"io/ioutil"

	"flag"

	log "github.com/Sirupsen/logrus"
)

var (
	config    *Config
	moduleLog = log.WithField("location", "config")
)

func GetConfig() *Config {
	if config == nil {
		moduleLog.Panic("Config has not initialized yet")
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
	Debug     bool              `json:"debug"`
	Port      int               `json:"port"`
	AccessLog bool              `json:"access_log"`
	Postgres  *PostgresSettings `json:"postgres"`
	Redis     *RedisSettings    `json:"redis"`
}

func getConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "conf", "", "Path to json config file.")
	flag.Parse()
	if configPath == "" {
		log.Panic("Cannot start without path to config")
	}
	return configPath
}
func Initialization() {
	path := getConfigPath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		moduleLog.WithField("config_path", path).Panicf("Cannot read config: %s", err)
	}
	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		moduleLog.Panicf("Cannot parse config: %s", err)
	}
}
