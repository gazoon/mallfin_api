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

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type Logo struct {
	Large string `json:"large"`
	Small string `json:"small"`
}
type SubwayStation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type WeekTime struct {
	Time time.Time `json:"time"`
	Day  int       `json:"day"`
}
type WorkPeriod struct {
	Open  WeekTime `json:"opening"`
	Close WeekTime `json:"closing"`
}
type Mall struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Phone         string         `json:"phone"`
	Address       string         `json:"address"`
	Logo          Logo           `json:"logo"`
	Location      Location       `json:"location"`
	SubwayStation *SubwayStation `json:"subway_station"`
}
type MallDetails struct {
	*Mall
	Site         string        `json:"site"`
	DayAndNight  bool          `json:"day_and_night"`
	WorkingHours []*WorkPeriod `json:"working_hours"`
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
	var subwayID *int
	var subwayName *string
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
	WHERE m.id = $1`, mallID).Scan(&mall.ID, &mall.Name, &mall.Phone, &mall.Address, &mall.Logo.Small, &mall.Logo.Large,
		&mall.Location.Lat, &mall.Location.Lon, &subwayID, &subwayName, &mall.Site, &mall.DayAndNight)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("mall", mallID).Panicf("Cannot get mall by ID: %s", err)
	}
	if subwayID != nil && subwayName != nil {
		mall.SubwayStation = &SubwayStation{Name: *subwayName, ID: *subwayID}
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
			err = rows.Scan(&period.Open.Day, &period.Open.Time, &period.Close.Day, &period.Close.Time)
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
