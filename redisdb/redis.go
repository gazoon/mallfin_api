package redisdb

import (
	"fmt"
	"mallfin_api/config"

	"sync"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

const NUMBER_OF_DATABASES = 16

var (
	db        *redis.Client
	moduleLog = log.WithField("location", "redis")
	once      sync.Once
)

func CreateNewDB(dbNumber int) *redis.Client {
	dbConf := config.Redis()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", dbConf.Host, dbConf.Port),
		Password: dbConf.Password,
		DB:       dbNumber,
	})
	err := client.Ping().Err()
	if err != nil {
		moduleLog.WithFields(log.Fields{"conf": dbConf, "db": dbNumber}).Panicf("Cannot connect to redis: %s", err)
	}
	return client
}

func GetClient() *redis.Client {
	if db == nil {
		moduleLog.Panic("Redis has not initialized yet.")
	}
	return db
}

func Initialization() {
	once.Do(func() {
		newDB := CreateNewDB(config.Redis().DB)
		db = newDB
	})
}

func Close() {
	if db != nil {
		db.Close()
	}
}
