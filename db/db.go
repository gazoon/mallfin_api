package db

import (
	"fmt"
	"mallfin_api/config"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-pg/pg"
)

var (
	db        *pg.DB
	moduleLog = log.WithField("location", "db")
	once      sync.Once
)

func CreateNewDB() *pg.DB {
	pgConf := config.Postgres()
	timeout := time.Second * time.Duration(pgConf.Timeout)
	db := pg.Connect(&pg.Options{
		Addr:         fmt.Sprintf("%s:%d", pgConf.Host, pgConf.Port),
		User:         pgConf.User,
		Password:     pgConf.Password,
		Database:     pgConf.DBName,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		PoolTimeout:  timeout,
		PoolSize:     pgConf.PoolSize,
		MaxRetries:   pgConf.Retries,
	})
	_, err := db.Exec(`SELECT 1`)
	if err != nil {
		moduleLog.WithFields(log.Fields{"conf": pgConf}).Panicf("Cannot connect to postgresql: %s", err)
	}
	return db
}

func Initialization() {
	once.Do(func() {
		newDB := CreateNewDB()
		db = newDB
	})
}

func GetClient() *pg.DB {
	if db == nil {
		moduleLog.Panic("Postgres has not initialized yet")
	}
	return db
}

func Close() {
	if db != nil {
		db.Close()
	}
}
