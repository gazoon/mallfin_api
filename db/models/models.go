package models

import (
	"database/sql"
	"mallfin_api/db"

	log "github.com/Sirupsen/logrus"
)

var (
	moduleLog = log.WithField("location", "models")
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type Mall struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	Site        *string  `json:"site"`
	DayAndNight bool     `json:"day_and_night"`
	Location    Location `json:"location"`
	CityID      int      `json:"-"`
}

func DeleteAllMalls() {
	conn := db.GetConnection()
	_, err := conn.Exec(`TRUNCATE mall CASCADE`)
	if err != nil {
		moduleLog.Panicf("Cannot delete malls: %s", err)
	}
}
func CreateMall(mall *Mall) *Mall {
	conn := db.GetConnection()
	err := conn.QueryRow(`
	INSERT INTO mall (name, site, address, day_and_night, city_id, location)
	VALUES ($1, $2, $3, $4, $5, ST_SETSRID(ST_POINT($6, $7), 4326))
	RETURNING id`, mall.Name, mall.Site, mall.Address, mall.DayAndNight, mall.CityID, mall.Location.Lat, mall.Location.Lon).Scan(&mall.ID)
	if err != nil {
		moduleLog.WithField("mall", mall).Panicf("Cannot craete mall: %s", err)
	}
	return mall

}
func GetMall(mallID int) *Mall {
	conn := db.GetConnection()
	mall := Mall{}
	err := conn.QueryRow(`
	SELECT
	  id,
	  name,
	  address,
	  site,
	  day_and_night,
	  ST_X(location),
	  ST_Y(location)
	FROM mall
	WHERE id = $1`, mallID).Scan(&mall.ID, &mall.Name, &mall.Address, &mall.Site, &mall.DayAndNight, &mall.Location.Lat, &mall.Location.Lon)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("mall", mallID).Panicf("Cannot get mall by ID: %s", err)
	}
	return &mall
}
