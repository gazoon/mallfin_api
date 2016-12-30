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

type PostgresSettings struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
type RedisSettings struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
type config struct {
	Debug    bool              `json:"debug"`
	Postgres *PostgresSettings `json:"postgres"`
	Redis    *RedisSettings    `json:"redis"`
}

func Initialization(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithFields(log.Fields{"location": "config", "config_path": path}).Fatalf("Cannot read config: %s", err)
	}
	conf = &config{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		log.WithFields(log.Fields{"location": "config"}).Fatalf("Cannot parse config: %s", err)
	}
}
