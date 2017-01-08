package models

import (
	"database/sql"
	"mallfin_api/db"

	"time"

	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
)

var moduleLog = log.WithField("location", "models")

type OrderBy struct {
	Column  string
	Reverse bool
}

func (o *OrderBy) String() string {
	var s string
	if o.Reverse {
		s = fmt.Sprintf("%s DESC", o.Column)
	} else {
		s = fmt.Sprintf("%s ASC", o.Column)
	}
	return s
}

type SortKeyToOrderBy map[string]*OrderBy

var (
	MALL_DEFAULT_ORDER_BY = &OrderBy{Column: "m.id", Reverse: false}
	SHOP_DEFAULT_ORDER_BY = &OrderBy{Column: "s.id", Reverse: false}
	MALL_SORT_KEYS        = SortKeyToOrderBy{
		"id":           MALL_DEFAULT_ORDER_BY,
		"-id":          {Column: "m.id", Reverse: true},
		"name":         {Column: "m.name", Reverse: false},
		"-name":        {Column: "m.name", Reverse: true},
		"shops_count":  {Column: "m.shops_count", Reverse: false},
		"-shops_count": {Column: "m.shops_count", Reverse: true},
	}
	SHOP_SORT_KEYS = SortKeyToOrderBy{
		"id":     SHOP_DEFAULT_ORDER_BY,
		"-id":    {Column: "s.id", Reverse: true},
		"name":   {Column: "s.name", Reverse: false},
		"-name":  {Column: "s.name", Reverse: true},
		"score":  {Column: "s.score", Reverse: false},
		"-score": {Column: "s.score", Reverse: true},
	}
)

func (sk SortKeyToOrderBy) FmtKeys() string {
	keys := make([]string, 0, len(sk))
	for key := range sk {
		keys = append(keys, key)
	}
	return strings.Join(keys, ", ")
}
func (sk SortKeyToOrderBy) CorrespondingOrderBy(sortKey *string) *OrderBy {
	orderBy := MALL_DEFAULT_ORDER_BY
	if sortKey != nil {
		if correspondOrderBy, ok := MALL_SORT_KEYS[*sortKey]; ok {
			orderBy = correspondOrderBy
		}
	}
	return orderBy
}

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

func ExistsQuery(query string, args ...interface{}) bool {
	var exists bool
	conn := db.GetConnection()
	err := conn.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		moduleLog.Panicf("Cannot check the existence: %s", err)
	}
	return exists
}
func IsShopExists(shopID int) bool {
	exists := ExistsQuery(`
	SELECT exists(
		SELECT *
		FROM shop
		WHERE id = $1)
	`, shopID)
	return exists
}
func IsMallExists(mallID int) bool {
	exists := ExistsQuery(`
	SELECT exists(
		SELECT *
		FROM mall
		WHERE id = $1)
	`, mallID)
	return exists
}
func IsCityExists(cityID int) bool {
	exists := ExistsQuery(`
	SELECT exists(
		SELECT *
		FROM city
		WHERE id = $1)
	`, cityID)
	return exists
}
func IsCategoryExists(categoryID int) bool {
	exists := ExistsQuery(`
	SELECT exists(
		SELECT *
		FROM category
		WHERE id = $1)
	`, categoryID)
	return exists
}
func IsSubwayStationExists(subwayStationID int) bool {
	exists := ExistsQuery(`
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
func GetMallWorkingHours(mallID int) []*WorkPeriod {
	locLog := moduleLog.WithField("mall", mallID)
	conn := db.GetConnection()
	rows, err := conn.Query(`
		SELECT
		  opening_day,
		  opening_time,
		  closing_day,
		  closing_time
		FROM mall_working_hours
		WHERE mall_id = $1
		`, mallID)
	if err != nil && err != sql.ErrNoRows {
		locLog.Panicf("Cannot get mall working hours: %s", err)
	}
	defer rows.Close()
	var workingHours []*WorkPeriod
	for rows.Next() {
		period := WorkPeriod{}
		err = rows.Scan(&period.OpenDay, &period.OpenTime, &period.CloseDay, &period.CloseTime)
		if err != nil {
			locLog.Panicf("Error during scaning working hours: %s", err)
		}
		workingHours = append(workingHours, &period)
	}
	err = rows.Err()
	if err != nil {
		locLog.Panicf("Error after scaning working hours: %s", err)
	}
	return workingHours
}
func GetMallDetails(mallID int) *Mall {
	conn := db.GetConnection()
	mall := Mall{}
	err := conn.QueryRow(`
	SELECT
	  m.id,
	  m.name,
	  m.phone,
	  m.logo_small,
	  m.logo_large,
	  ST_X(m.location) location_lat,
	  ST_Y(m.location) location_lon,
	  m.shops_count,
	  m.address,
	  m.site,
	  m.day_and_night,
	  m.subway_station_id,
	  ss.name          subway_station_name
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.id
	WHERE m.id = $1
	`, mallID).Scan(&mall.ID, &mall.Name, &mall.Phone, &mall.LogoSmall, &mall.LogoLarge, &mall.LocationLat, &mall.LocationLon, &mall.ShopsCount,
		&mall.Address, &mall.Site, &mall.DayAndNight, &mall.SubwayID, &mall.SubwayName)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("mall", mallID).Panicf("Cannot get mall by ID: %s", err)
	}
	if !mall.DayAndNight {
		mall.WorkingHours = GetMallWorkingHours(mall.ID)
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
	orderBy := MALL_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		malls = MallsQuery(fmt.Sprintf(`
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
		WHERE m.city_id = $3
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`, orderBy), limit, offset, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		WHERE m.city_id = $1
		`, *cityID)
	} else {
		malls = MallsQuery(fmt.Sprintf(`
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
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`, orderBy), limit, offset)
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
		return nil, 0
	}
	orderBy := MALL_SORT_KEYS.CorrespondingOrderBy(sortKey)
	malls := MallsQuery(fmt.Sprintf(`
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
	WHERE m.id = ANY($3)
	ORDER BY %s
	LIMIT $1
	OFFSET $2
	`, orderBy), limit, offset, pq.Array(mallIDs))
	totalCount := countQuery(`
	SELECT
	  count(*) total_count
	FROM mall m
	WHERE m.id = ANY($1)
	`, pq.Array(mallIDs))
	return malls, totalCount
}
func GetMallsBySubwayStation(subwayStationID int, sortKey *string, limit, offset *uint) ([]*Mall, int) {
	orderBy := MALL_SORT_KEYS.CorrespondingOrderBy(sortKey)
	malls := MallsQuery(fmt.Sprintf(`
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
	WHERE ss.id = $3
	ORDER BY %s
	LIMIT $1
	OFFSET $2
	`, orderBy), limit, offset, subwayStationID)
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
	orderBy := MALL_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		malls = MallsQuery(fmt.Sprintf(`
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
		WHERE ms.shop_id = $3 AND m.city_id = $4
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`, orderBy), limit, offset, shopID, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = $1 AND m.city_id = $2
		`, shopID, *cityID)
	} else {
		malls = MallsQuery(fmt.Sprintf(`
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
		WHERE ms.shop_id = $3
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`, orderBy), limit, offset, shopID)
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
	orderBy := MALL_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		malls = MallsQuery(fmt.Sprintf(`
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
		  JOIN (SELECT DISTINCT ON (mall_id) mall_id
				FROM mall_name
				WHERE name ILIKE '%%' || $3 || '%%') mn ON m.id = mn.mall_id
		WHERE m.city_id = $4
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`, orderBy), limit, offset, name, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		  JOIN (SELECT DISTINCT ON (mall_id) mall_id
				FROM mall_name
				WHERE name ILIKE '%' || $1 || '%') mn ON m.id = mn.mall_id
		WHERE m.city_id = $2
		`, name, *cityID)
	} else {
		malls = MallsQuery(fmt.Sprintf(`
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
		  JOIN (SELECT DISTINCT ON (mall_id) mall_id
				FROM mall_name
				WHERE name ILIKE '%%' || $3 || '%%') mn ON m.id = mn.mall_id
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`, orderBy), limit, offset, name)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		  JOIN (SELECT DISTINCT ON (mall_id) mall_id
				FROM mall_name
				WHERE name ILIKE '%' || $1 || '%') mn ON m.id = mn.mall_id
		`, name)
	}
	return malls, totalCount
}
func MallsQuery(query string, args ...interface{}) []*Mall {
	conn := db.GetConnection()
	rows, err := conn.Query(query, args...)
	if err != nil {
		moduleLog.Panicf("Cannot get malls rows: %s", err)
	}
	defer rows.Close()
	var malls []*Mall
	for rows.Next() {
		m := Mall{}
		err = rows.Scan(&m.ID, &m.Name, &m.Phone, &m.LogoSmall, &m.LogoLarge, &m.LocationLat, &m.LocationLon, &m.ShopsCount)
		if err != nil {
			moduleLog.Panicf("Error during scaning mall row: %s", err)
		}
		malls = append(malls, &m)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning malls rows: %s", err)
	}
	return malls
}
