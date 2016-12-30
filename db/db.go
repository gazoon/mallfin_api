package db

import (
	"database/sql"
	"fmt"
	"mallfin_api/config"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func Initialization() {
	conf := config.Postgres()
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d", conf.Name, conf.User, conf.Password, conf.Host, conf.Port))
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		log.WithFields(log.Fields{"location": "postgres", "conf": conf}).Panicf("Cannot connect to postgresql: %s", err)
	}
}
func GetConnection() *sql.DB {
	if db == nil {
		log.WithField("location", "postgres").Panic("Postgres has not initialized yet")
	}
	return db
}
func Close() {
	db.Close()
}
