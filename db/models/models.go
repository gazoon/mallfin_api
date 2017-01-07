package models

import (
	"database/sql"
	"mallfin_api/db"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
)

var (
	moduleLog = log.WithField("location", "models")
)

const (
	NAME_MALL_SORT_KEY        = "name"
	SHOPS_COUNT_MALL_SORT_KEY = "shops_count"
	NAME_SHOP_SORT_KEY        = "name"
	SCORE_SHOP_SORT_KEY       = "score"
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
	// details:
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
func GetMallDetails(mallID int) *Mall {
	conn := db.GetConnection()
	mall := Mall{}
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
func countQuery(query string, args ...interface{}) int {
	var count int
	conn := db.GetConnection()
	err := conn.QueryRow(query, args...).Scan(&count)
	if err != nil {
		moduleLog.Panicf("Cannot do count query: %s", err)
	}
	return count
}
func GetMalls(cityID *int, sortKey *string, limit, offset *uint) ([]*Mall, int) {
	var malls []*Mall
	var totalCount int
	if cityID != nil {
		malls = mallsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location)  location_lat,
		  ST_Y(m.location)  location_lon,
		  m.shops_count
		FROM mall m
		WHERE m.city_id = $1
		`, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		WHERE m.city_id = $1
		`, *cityID)
	} else {
		malls = mallsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location)  location_lat,
		  ST_Y(m.location)  location_lon,
		  m.shops_count
		FROM mall m
		`)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		`)
	}
	return malls, totalCount
}
func GetMallsByIds(mallIDs []int, sortKey *string, limit, offset *uint) ([]*Mall, int) {
	if len(mallIDs) == 0 {
		return []*Mall{}, 0
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
	  m.shops_count
	FROM mall m
	WHERE m.id = ANY($1)
	`, pq.Array(mallIDs))
	totalCount := countQuery(`
	SELECT
	  count(*) total_count
	FROM mall m
	WHERE m.id = ANY($1)
	`, pq.Array(mallIDs))
	return malls, totalCount
}
func GetMallsBySubwayStation(subwayStationID int, sortKey *string, limit, offset *uint) ([]*Mall, int) {
	malls := mallsQuery(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location)  location_lat,
	  ST_Y(m.location)  location_lon,
	  m.shops_count
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.id
	WHERE ss.id = $1
	`, subwayStationID)
	totalCount := countQuery(`
	SELECT
	  count(*) total_count
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.id
	WHERE ss.id = $1
	`, subwayStationID)
	return malls, totalCount
}
func GetMallsByShop(shopID int, cityID *int, sortKey *string, limit, offset *uint) ([]*Mall, int) {
	var malls []*Mall
	var totalCount int
	if cityID != nil {
		malls = mallsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location) location_lat,
		  ST_Y(m.location) location_lon,
		  m.shops_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = $1 AND m.city_id = $2
		`, shopID, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = $1 AND m.city_id = $2
		`, shopID, *cityID)
	} else {
		malls = mallsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location) location_lat,
		  ST_Y(m.location) location_lon,
		  m.shops_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = $1
		`, shopID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = $1
		`, shopID)
	}
	return malls, totalCount
}
func GetMallsByName(name string, cityID *int, sortKey *string, limit, offset *uint) ([]*Mall, int) {
	var malls []*Mall
	var totalCount int
	if cityID != nil {
		malls = mallsQuery(`
		SELECT DISTINCT ON (m.id)
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location) location_lat,
		  ST_Y(m.location) location_lon,
		  m.shops_count
		FROM mall m
		  JOIN mall_name mn ON m.id = mn.mall_id
		WHERE mn.name ILIKE '%' || $1 || '%' AND m.city_id = $2
		`, name, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(DISTINCT m.id) total_count
		FROM mall m
		  JOIN mall_name mn ON m.id = mn.mall_id
		WHERE mn.name ILIKE '%' || $1 || '%' AND m.city_id = $2
		`, name, *cityID)
	} else {
		malls = mallsQuery(`
		SELECT DISTINCT ON (m.id)
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location) location_lat,
		  ST_Y(m.location) location_lon,
		  m.shops_count
		FROM mall m
		  JOIN mall_name mn ON m.id = mn.mall_id
		WHERE mn.name ILIKE '%' || $1 || '%'
		`, name)
		totalCount = countQuery(`
		SELECT
		  count(DISTINCT m.id) total_count
		FROM mall m
		  JOIN mall_name mn ON m.id = mn.mall_id
		WHERE mn.name ILIKE '%' || $1 || '%'
		`, name)
	}
	return malls, totalCount
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
		malls = append(malls, &m)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning malls rows: %s", err)
	}
	return malls
}
