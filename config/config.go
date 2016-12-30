package config

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

var (
	conf *config
)

func checkInitialization() {
	if conf == nil {
		log.WithField("location", "config").Panic("Config has not initialized yet")
	}
}
func Postgres() *PostgresSettings {
	checkInitialization()
	return conf.Postgres
}
func Redis() *RedisSettings {
	checkInitialization()
	return conf.Redis
}
func Debug() bool {
	checkInitialization()
	return conf.Debug
}
func Port() int {
	checkInitialization()
	return conf.Port
}

type PostgresSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
type RedisSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
type config struct {
	Debug    bool              `json:"debug"`
	Port     int               `json:"port"`
	Postgres *PostgresSettings `json:"postgres"`
	Redis    *RedisSettings    `json:"redis"`
}

func Initialization(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithFields(log.Fields{"location": "config", "config_path": path}).Panicf("Cannot read config: %s", err)
	}
	conf = &config{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		log.WithFields(log.Fields{"location": "config"}).Panicf("Cannot parse config: %s", err)
	}
}
