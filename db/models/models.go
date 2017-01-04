package models

import (
	"database/sql"
	"mallfin_api/db"

	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	moduleLog = log.WithField("location", "models")
)

type WorkPeriod struct {
	OpenTime  time.Time
	OpenDay   int
	CloseTime time.Time
	CloseDay  int
}
type Mall struct {
	ID          int
	Name        string
	Phone       string
	Address     string
	LogoLarge   string
	LogoSmall   string
	LocationLat float64
	LocationLon float64
	SubwayID    *int
	SubwayName  *string
}
type MallDetails struct {
	*Mall
	Site         string
	DayAndNight  bool
	WorkingHours []*WorkPeriod
}

func DeleteAllMalls() {
	conn := db.GetConnection()
	_, err := conn.Exec(`TRUNCATE mall CASCADE`)
	if err != nil {
		moduleLog.Panicf("Cannot delete malls: %s", err)
	}
}

//func CreateMall(mall *Mall) *Mall {
//	conn := db.GetConnection()
//	err := conn.QueryRow(`
//	INSERT INTO mall (name, site, address, day_and_night, city_id, location)
//	VALUES ($1, $2, $3, $4, $5, ST_SETSRID(ST_POINT($6, $7), 4326))
//	RETURNING id`, mall.Name, mall.Site, mall.Address, mall.DayAndNight, mall.CityID, mall.Location.Lat, mall.Location.Lon).Scan(&mall.ID)
//	if err != nil {
//		moduleLog.WithField("mall", mall).Panicf("Cannot craete mall: %s", err)
//	}
//	return mall
//
//}
func GetMallDetails(mallID int) *MallDetails {
	conn := db.GetConnection()
	mall := MallDetails{Mall: new(Mall)}
	err := conn.QueryRow(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.address,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location),
	  ST_Y(m.location),
	  ss.id,
	  ss.name,
	  m.site,
	  m.day_and_night
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.id
	WHERE m.id = $1`, mallID).Scan(&mall.ID, &mall.Name, &mall.Phone, &mall.Address, &mall.LogoSmall, &mall.LogoLarge,
		&mall.LocationLat, &mall.LocationLon, &mall.SubwayID, &mall.SubwayName, &mall.Site, &mall.DayAndNight)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("mall", mallID).Panicf("Cannot get mall by ID: %s", err)
	}
	if !mall.DayAndNight {
		rows, err := conn.Query(`
		SELECT
		  opening_day,
		  opening_time,
		  closing_day,
		  closing_time
		FROM mall_working_hours
		WHERE mall_id = $1`, mall.ID)
		if err != nil && err != sql.ErrNoRows {
			moduleLog.WithField("mall", mall.ID).Panicf("Cannot get mall working hours: %s", err)
		}
		defer rows.Close()
		for rows.Next() {
			period := WorkPeriod{}
			err = rows.Scan(&period.OpenDay, &period.OpenTime, &period.CloseDay, &period.CloseTime)
			if err != nil {
				moduleLog.WithField("mall", mall.ID).Panicf("Error during scaning working hours: %s", err)
			}
			mall.WorkingHours = append(mall.WorkingHours, &period)
		}
		err = rows.Err()
		if err != nil {
			moduleLog.WithField("mall", mall.ID).Panicf("Error after scaning working hours: %s", err)
		}
	}
	return &mall
}

//func GetMallsByIds(mallIDs []int) []*Mall {
//	malls := []*Mall{}
//	if len(mallIDs) == 0 {
//		return malls
//	}
//	conn:=db.GetConnection()
//	rows,err:=conn.Query(`
//	SELECT
//	  m.id,
//	  m.name,
//	  m.phone,
//	  m.address,
//	  m.logo_small,
//	  m.logo_large,
//	  ST_X(m.location),
//	  ST_Y(m.location),
//	  ss.id,
//	  ss.name
//	FROM mall m LEFT JOIN subway_station ss ON m.subway_station_id = ss.id
//	WHERE m.id IN $1`,mallIDs)
//	if err!= nil {
//		moduleLog.Panicf("Cannot get malls by ids: %s",err)
//	}
//	defer rows.Close()
//	for rows.Next() {
//		mall:=Mall{}
//		err = rows.Scan(&mall.ID,mall.Name)
//		if err != nil {
//			moduleLog.WithField("mall", mall.ID).Panicf("Error during scaning working hours: %s", err)
//		}
//		mall.WorkingHours = append(mall.WorkingHours, &period)
//	}
//	err = rows.Err()
//	if err != nil {
//		moduleLog.WithField("mall", mall.ID).Panicf("Error after scaning working hours: %s", err)
//	}
//
//	return malls
//}
