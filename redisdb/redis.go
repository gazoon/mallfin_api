package redisdb

import (
	"fmt"
	"mallfin_api/config"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

const NUMBER_OF_DATABASES = 16

var (
	db        *redis.Client
	moduleLog = log.WithField("location", "redis")
)

func newDBConnection(dbNumber int) *redis.Client {
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
func closeCurrentConnection() {
	if db != nil {
		db.Close()
	}
}
func setNewConnection(conn *redis.Client) {
	closeCurrentConnection()
	db = conn
}
func GetConnection() *redis.Client {
	if db == nil {
		moduleLog.Panic("Redis has not initialized yet.")
	}
	return db
}

func Initialization() {
	conn := newDBConnection(config.Redis().DB)
	setNewConnection(conn)
}
func InitializationForTests() {
	realDBNumber := config.Redis().DB
	testDBNumber := (realDBNumber + 1) % NUMBER_OF_DATABASES
	conn := newDBConnection(testDBNumber)
	setNewConnection(conn)
}
func Close() {
	closeCurrentConnection()
}
