package db

import (
	"github.com/go-pg/pg"

	"fmt"
	"mallfin_api/config"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	db        *pg.DB
	moduleLog = log.WithField("location", "postgres")
)

func createNewDB() *pg.DB {
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
func closeCurrentDB() {
	if db != nil {
		db.Close()
	}
}
func setNewDB(newDB *pg.DB) {
	closeCurrentDB()
	db = newDB
}

func Initialization() {
	db := createNewDB()
	setNewDB(db)
}
func GetClient() *pg.DB {
	if db == nil {
		moduleLog.Panic("Postgres has not initialized yet")
	}
	return db
}
func Close() {
	closeCurrentDB()
}
