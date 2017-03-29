package redisdb

import (
	"fmt"
	"mallfin_api/config"

	"sync"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
	"mallfin_api/logging"
)

const NumberOfDatabases = 16

var (
	db     *redis.Client
	logger = logging.WithPackage("redis")
	once   sync.Once
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
		logger.WithFields(log.Fields{"conf": dbConf, "db": dbNumber}).Panicf("Cannot connect to redis: %s", err)
	}
	return client
}

func GetClient() *redis.Client {
	if db == nil {
		logger.Panic("Redis has not initialized yet.")
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
