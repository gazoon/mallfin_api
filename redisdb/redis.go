package redisdb

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
	"mallfin_api/config"
)

var (
	db *redis.Client
)

func NewDBConnection() *redis.Client {
	conf := config.Redis()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	err := client.Ping().Err()
	if err != nil {
		log.WithFields(log.Fields{"location": "redis", "conf": conf}).Panicf("Cannot connect to redis: %s", err)
	}
	return client
}
func GetConnection() *redis.Client {
	if db == nil {
		log.WithField("location", "redis").Panic("Redis has not initialized yet.")
	}
	return db
}
func Initialization() {
	db = NewDBConnection()
}
func Close() {
	db.Close()
}
