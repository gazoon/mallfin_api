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
	LogoLarge   string
	LogoSmall   string
	LocationLat float64
	LocationLon float64
	ShopsCount  int
}
type MallDetails struct {
	*Mall
	Address      string
	Site         string
	DayAndNight  bool
	SubwayID     *int
	SubwayName   *string
	WorkingHours []*WorkPeriod
}

func existsQuery(query string, args ...interface{}) bool {
	var exists bool
	conn := db.GetConnection()
	err := conn.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		moduleLog.Panicf("Cannot check the existence: %s", err)
	}
	return exists
}
func IsShopExists(shopID int) bool {
	exists := existsQuery(`
	SELECT exists(
		SELECT *
		FROM shop
		WHERE id = $1)
	`, shopID)
	return exists
}
func IsMallExists(mallID int) bool {
	exists := existsQuery(`
	SELECT exists(
		SELECT *
		FROM mall
		WHERE id = $1)
	`, mallID)
	return exists
}
func IsCityExists(cityID int) bool {
	exists := existsQuery(`
	SELECT exists(
		SELECT *
		FROM city
		WHERE id = $1)
	`, cityID)
	return exists
}
func IsCategoryExists(categoryID int) bool {
	exists := existsQuery(`
	SELECT exists(
		SELECT *
		FROM category
		WHERE id = $1)
	`, categoryID)
	return exists
}
func IsSubwayStationExists(subwayStationID int) bool {
	exists := existsQuery(`
	SELECT exists(
		SELECT *
		FROM subway_station
		WHERE id = $1)
	`, subwayStationID)
	return exists
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

func GetMallsByIds(mallIDs []int) []*Mall {
	if len(mallIDs) == 0 {
		return []*Mall{}
	}
	malls := mallsQuery(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location)  location_lat,
	  ST_Y(m.location)  location_lon,
	  count(ms.shop_id) shops_count
	FROM mall m
	  JOIN mall_shop ms ON m.id = ms.mall_id
	WHERE m.id IN $1
	GROUP BY m.id
	`, mallIDs)
	return malls
}
func GetMallsBySubwayStation(subwayStationID int) []*Mall {
	malls := mallsQuery(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location)  location_lat,
	  ST_Y(m.location)  location_lon,
	  count(ms.shop_id) shops_count
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.id
	  JOIN mall_shop ms ON m.id = ms.mall_id
	WHERE ss.id = $1
	GROUP BY m.id
	`, subwayStationID)
	return malls
}
func GetMallsByShop(shopID int) []*Mall {
	malls := mallsQuery(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location)  location_lat,
	  ST_Y(m.location)  location_lon,
	  count(ms.shop_id) shops_count
	FROM mall m
	  JOIN (SELECT mall_id
			FROM mall_shop
			WHERE shop_id = $1) q ON m.id = q.mall_id
	  JOIN mall_shop ms ON m.id = ms.mall_id
	GROUP BY m.id
	`, shopID)
	return malls
}
func GetMallsByShopAndCity(shopID, cityID int) []*Mall {
	malls := mallsQuery(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location)  location_lat,
	  ST_Y(m.location)  location_lon,
	  count(ms.shop_id) shops_count
	FROM mall m
	  JOIN (SELECT mall_id
			FROM mall_shop
			WHERE shop_id = $1) q ON m.id = q.mall_id
	  JOIN mall_shop ms ON m.id = ms.mall_id
	WHERE m.city_id = $2
	GROUP BY m.id
	`, shopID, cityID)
	return malls
}
func GetMallsByName(name string) []*Mall {
	malls := mallsQuery(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.address,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location)  location_lat,
	  ST_Y(m.location)  location_lon,
	  count(ms.shop_id) shops_count
	FROM mall m
	  JOIN (SELECT DISTINCT ON (mn.mall_id) mn.mall_id
			FROM mall_name mn
			WHERE mn.name ILIKE '%' || $1 || '%'
			ORDER BY mn.mall_id) mn ON m.id = mn.mall_id
	  JOIN mall_shop ms ON m.id = ms.mall_id
	GROUP BY m.id
	`, name)
	return malls
}
func GetMallsByNameAndCity(name string, cityID int) []*Mall {
	malls := mallsQuery(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.address,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location)  location_lat,
	  ST_Y(m.location)  location_lon,
	  count(ms.shop_id) shops_count
	FROM mall m
	  JOIN (SELECT DISTINCT ON (mn.mall_id) mn.mall_id
			FROM mall_name mn
			WHERE mn.name ILIKE '%' || $1 || '%'
			ORDER BY mn.mall_id) mn ON m.id = mn.mall_id
	  JOIN mall_shop ms ON m.id = ms.mall_id
	WHERE m.city_id = $2
	GROUP BY m.id
	`, name, cityID)
	return malls
}
func mallsQuery(query string, args ...interface{}) []*Mall {
	conn := db.GetConnection()
	rows, err := conn.Query(query, args...)
	if err != nil {
		moduleLog.Panicf("Cannot get malls rows: %s", err)
	}
	defer rows.Close()
	malls := []*Mall{}
	for rows.Next() {
		m := Mall{}
		err = rows.Scan(&m.ID, &m.Name, &m.Phone, &m.LogoSmall, &m.LogoLarge, &m.LocationLat, &m.LocationLon, &m.ShopsCount)
		if err != nil {
			moduleLog.Panicf("Error during mall row: %s", err)
		}
		malls = append(malls, m)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning malls rows: %s", err)
	}
	return malls
}
