package models

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"mallfin_api/db"
)

var (
	moduleLog = log.WithField("location", "models")
)

type Mall struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"string"`
	DayAndNight bool   `json:"day_and_night"`
}

func GetMall(mallID int) *Mall {
	conn := db.GetConnection()
	mall := Mall{}
	err := conn.QueryRow(`
	SELECT
	  id,
	  name,
	  address,
	  day_and_night
	FROM mall
	WHERE id = $1`, mallID).Scan(&mall.ID, &mall.Name, &mall.Address, &mall.DayAndNight)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("mall", mallID).Panicf("Cannot get mall by ID: %s", err)
	}
	return &mall
}
