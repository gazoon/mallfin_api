package models

import (
	"database/sql"
	"mallfin_api/db"

	"time"

	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/pq"
)

var moduleLog = log.WithField("location", "models")

type OrderBy struct {
	Column  string
	Reverse bool
}

func (o *OrderBy) String() string {
	return o.ToSql()
}
func (o *OrderBy) ToSql() string {
	var s string
	if o.Reverse {
		s = fmt.Sprintf("%s DESC", o.Column)
	} else {
		s = fmt.Sprintf("%s ASC", o.Column)
	}
	return s

}
func (o *OrderBy) Compile(query string) string {
	return fmt.Sprintf(query, o.ToSql())
}

type SortKeyToOrderBy struct {
	dict           map[string]*OrderBy
	defaultOrderBy *OrderBy
}

var (
	MALL_DEFAULT_ORDER_BY     = &OrderBy{Column: "m.id", Reverse: false}
	SHOP_DEFAULT_ORDER_BY     = &OrderBy{Column: "s.id", Reverse: false}
	CATEGORY_DEFAULT_ORDER_BY = &OrderBy{Column: "c.id", Reverse: false}
	CITY_DEFAULT_ORDER_BY     = &OrderBy{Column: "c.id", Reverse: false}

	MALLS_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           MALL_DEFAULT_ORDER_BY,
			"-id":          {Column: "m.id", Reverse: true},
			"name":         {Column: "m.name", Reverse: false},
			"-name":        {Column: "m.name", Reverse: true},
			"shops_count":  {Column: "m.shops_count", Reverse: false},
			"-shops_count": {Column: "m.shops_count", Reverse: true},
		},
		defaultOrderBy: MALL_DEFAULT_ORDER_BY,
	}
	SHOPS_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           SHOP_DEFAULT_ORDER_BY,
			"-id":          {Column: "s.id", Reverse: true},
			"name":         {Column: "s.name", Reverse: false},
			"-name":        {Column: "s.name", Reverse: true},
			"score":        {Column: "s.score", Reverse: false},
			"-score":       {Column: "s.score", Reverse: true},
			"malls_count":  {Column: "s.malls_count", Reverse: false},
			"-malls_count": {Column: "s.malls_count", Reverse: true},
		},
		defaultOrderBy: SHOP_DEFAULT_ORDER_BY,
	}
	CATEGORIES_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           CATEGORY_DEFAULT_ORDER_BY,
			"-id":          {Column: "c.id", Reverse: true},
			"name":         {Column: "c.name", Reverse: false},
			"-name":        {Column: "c.name", Reverse: true},
			"shops_count":  {Column: "c.shops_count", Reverse: false},
			"-shops_count": {Column: "c.shops_count", Reverse: true},
		},
		defaultOrderBy: CATEGORY_DEFAULT_ORDER_BY,
	}
	CITIES_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":    CITY_DEFAULT_ORDER_BY,
			"-id":   {Column: "c.id", Reverse: true},
			"name":  {Column: "c.name", Reverse: false},
			"-name": {Column: "c.name", Reverse: true},
		},
		defaultOrderBy: CITY_DEFAULT_ORDER_BY,
	}
)

func (sk *SortKeyToOrderBy) FmtKeys() string {
	keys := make([]string, 0, len(sk.dict))
	for key := range sk.dict {
		keys = append(keys, key)
	}
	return strings.Join(keys, ", ")
}
func (sk *SortKeyToOrderBy) IsValid(sortKey *string) bool {
	if sortKey != nil {
		if _, ok := sk.dict[*sortKey]; !ok {
			return false
		}
	}
	return true
}
func (sk *SortKeyToOrderBy) CorrespondingOrderBy(sortKey *string) *OrderBy {
	orderBy := sk.defaultOrderBy
	if sortKey != nil {
		if correspondOrderBy, ok := sk.dict[*sortKey]; ok {
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
type Location struct {
	Lat float64
	Lon float64
}
type Mall struct {
	ID         int
	Name       string
	Phone      string
	LogoLarge  string
	LogoSmall  string
	Location   Location
	ShopsCount int
	Address    string
	//Details
	Site         string
	DayAndNight  bool
	SubwayID     *int
	SubwayName   *string
	WorkingHours []*WorkPeriod
}

type Shop struct {
	ID         int
	Name       string
	LogoLarge  string
	LogoSmall  string
	Score      int
	MallsCount int
	//Details
	Phone       string
	Site        string
	NearestMall *int
}

type Category struct {
	ID         int
	Name       string
	LogoLarge  string
	LogoSmall  string
	ShopsCount int
}

type City struct {
	ID   int
	Name string
}
type ShopsInMalls map[int][]int

type SearchResult struct {
	Mall     *Mall
	ShopIDs  []int
	Distance *float64
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

//func DeleteAllMalls() {
//	conn := db.GetConnection()
//	_, err := conn.Exec(`TRUNCATE mall CASCADE`)
//	if err != nil {
//		moduleLog.Panicf("Cannot delete malls: %s", err)
//	}
//}

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
func searchResultsQuery(query string, args ...interface{}) []*SearchResult {
	conn := db.GetConnection()
	rows, err := conn.Query(query, args...)
	if err != nil {
		moduleLog.Panicf("Cannot get search results rows: %s", err)
	}
	defer rows.Close()
	var searchResults []*SearchResult
	for rows.Next() {
		sr := SearchResult{Mall: &Mall{}}
		err = rows.Scan(&sr.Mall.ID, &sr.Mall.Name, &sr.Mall.Phone, &sr.Mall.LogoSmall, &sr.Mall.LogoLarge, &sr.Mall.Location.Lat, &sr.Mall.Location.Lon, &sr.Mall.ShopsCount, pq.Array(&sr.ShopIDs), &sr.Distance)
		if err != nil {
			moduleLog.Panicf("Error during scaning search result row: %s", err)
		}
		searchResults = append(searchResults, &sr)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning search results rows: %s", err)
	}
	return searchResults

}
func GetSearchResults(shopIDs []int, cityID *int, limit, offset *uint) ([]*SearchResult, int) {
	if len(shopIDs) == 0 {
		return nil, 0
	}
	var searchResults []*SearchResult
	var totalCount int
	shopIDsArray := pq.Array(shopIDs)
	if cityID != nil {
		searchResults = searchResultsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location)      location_lat,
		  ST_Y(m.location)      location_lon,
		  m.shops_count         mall_shops_count,
		  array_agg(ms.shop_id) shops,
		  NULL                  distance
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($3) AND m.city_id = $4
		GROUP BY m.id
		ORDER BY count(ms.shop_id) DESC, mall_shops_count DESC
		LIMIT $1
		OFFSET $2
		`, limit, offset, shopIDsArray, *cityID)
		totalCount = countQuery(`
		SELECT count(DISTINCT m.id) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($1) AND m.city_id = $2
		`, shopIDsArray, *cityID)
	} else {
		searchResults = searchResultsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location)      location_lat,
		  ST_Y(m.location)      location_lon,
		  m.shops_count         mall_shops_count,
		  array_agg(ms.shop_id) shops,
		  NULL                  distance
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($3)
		GROUP BY m.id
		ORDER BY count(ms.shop_id) DESC, mall_shops_count DESC
		LIMIT $1
		OFFSET $2
		`, limit, offset, shopIDsArray)
		totalCount = countQuery(`
		SELECT count(DISTINCT m.id) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($1)
		`, shopIDsArray)
	}
	return searchResults, totalCount
}
func GetSearchResultsWithDistance(shopIDs []int, location *Location, cityID *int, limit, offset *uint) ([]*SearchResult, int) {
	if len(shopIDs) == 0 || location == nil {
		return nil, 0
	}
	var searchResults []*SearchResult
	var totalCount int
	shopIDsArray := pq.Array(shopIDs)
	if cityID != nil {
		searchResults = searchResultsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location)      location_lat,
		  ST_Y(m.location)      location_lon,
		  m.shops_count         mall_shops_count,
		  array_agg(ms.shop_id) shops,
		  st_distance(
			  st_transform(m.location, 26986),
			  st_transform(st_setsrid(st_point($4, $5), 4326), 26986)
		  )                     distance
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($3) AND m.city_id = $6
		GROUP BY m.id
		ORDER BY count(ms.shop_id) DESC, distance ASC
		LIMIT $1
		OFFSET $2
		`, limit, offset, shopIDsArray, location.Lat, location.Lon, *cityID)
		totalCount = countQuery(`
		SELECT count(DISTINCT m.id) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($1) AND m.city_id = $2
		`, shopIDsArray, *cityID)
	} else {
		searchResults = searchResultsQuery(`
		SELECT
		  m.id,
		  m.name,
		  m.phone,
		  m.logo_small,
		  m.logo_large,
		  ST_X(m.location)      location_lat,
		  ST_Y(m.location)      location_lon,
		  m.shops_count         mall_shops_count,
		  array_agg(ms.shop_id) shops,
		  st_distance(
			  st_transform(m.location, 26986),
			  st_transform(st_setsrid(st_point($4, $5), 4326), 26986)
		  )                     distance
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($3)
		GROUP BY m.id
		ORDER BY count(ms.shop_id) DESC, distance ASC
		LIMIT $1
		OFFSET $2
		`, limit, offset, shopIDsArray, location.Lat, location.Lon)
		totalCount = countQuery(`
		SELECT count(DISTINCT m.id) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($1)
		`, shopIDsArray)
	}
	return searchResults, totalCount
}
func GetShopsInMalls(mallIDs, shopIDs []int) ShopsInMalls {
	mallsShops := ShopsInMalls{}
	for _, mallID := range mallIDs {
		mallsShops[mallID] = nil
	}
	if len(mallIDs) == 0 || len(shopIDs) == 0 {
		return mallsShops
	}
	conn := db.GetConnection()
	rows, err := conn.Query(`
	SELECT
	  mall_id,
	  shop_id
	FROM mall_shop
	WHERE mall_id = ANY ($1) AND shop_id = ANY ($2)
	`, pq.Array(mallIDs), pq.Array(shopIDs))
	if err != nil && err != sql.ErrNoRows {
		moduleLog.Panicf("Cannot get shops in malls occurrence: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var mallID, shopID int
		err = rows.Scan(&mallID, &shopID)
		if err != nil {
			moduleLog.Panicf("Error during scaning shop in mall row: %s", err)
		}
		mallsShops[mallID] = append(mallsShops[mallID], shopID)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning shops in malls: %s", err)
	}
	return mallsShops
}
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
func mallQuery(query string, args ...interface{}) *Mall {
	conn := db.GetConnection()
	mall := Mall{}
	err := conn.QueryRow(query, args...).Scan(&mall.ID, &mall.Name, &mall.Phone, &mall.LogoSmall, &mall.LogoLarge, &mall.Location.Lat, &mall.Location.Lon, &mall.ShopsCount,
		&mall.Address, &mall.Site, &mall.DayAndNight, &mall.SubwayID, &mall.SubwayName)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.Panicf("Cannot get mall: %s", err)
	}
	if !mall.DayAndNight {
		mall.WorkingHours = GetMallWorkingHours(mall.ID)
	}
	return &mall
}
func GetMallDetails(mallID int) *Mall {
	mall := mallQuery(`
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
	`, mallID)
	return mall
}
func GetMallByLocation(location *Location) *Mall {
	if location == nil {
		return nil
	}
	mall := mallQuery(`
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
	WHERE st_dwithin(st_transform(m.location, 26986), st_transform(ST_Setsrid(st_point($1, $2), 4326), 26986), m.radius)
	ORDER BY m.location <-> ST_SetSRID(ST_Point($1, $2), 4326)
	LIMIT 1
	`, location.Lat, location.Lon)
	return mall

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
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		malls = mallsQuery(orderBy.Compile(`
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
		`), limit, offset, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		WHERE m.city_id = $1
		`, *cityID)
	} else {
		malls = mallsQuery(orderBy.Compile(`
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
		`), limit, offset)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		`)
	}
	return malls, totalCount
}
func GetMallsByIds(mallIDs []int) ([]*Mall, int) {
	if len(mallIDs) == 0 {
		return nil, 0
	}
	mallIDsArray := pq.Array(mallIDs)
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
	`, mallIDsArray)
	totalCount := countQuery(`
	SELECT
	  count(*) total_count
	FROM mall m
	WHERE m.id = ANY($1)
	`, mallIDsArray)
	return malls, totalCount
}
func GetMallsBySubwayStation(subwayStationID int, sortKey *string, limit, offset *uint) ([]*Mall, int) {
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	malls := mallsQuery(orderBy.Compile(`
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
	`), limit, offset, subwayStationID)
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
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		malls = mallsQuery(orderBy.Compile(`
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
		`), limit, offset, shopID, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM mall m
		  JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = $1 AND m.city_id = $2
		`, shopID, *cityID)
	} else {
		malls = mallsQuery(orderBy.Compile(`
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
		`), limit, offset, shopID)
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
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		malls = mallsQuery(orderBy.Compile(`
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
		`), limit, offset, name, *cityID)
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
		malls = mallsQuery(orderBy.Compile(`
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
		`), limit, offset, name)
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
func mallsQuery(query string, args ...interface{}) []*Mall {
	conn := db.GetConnection()
	rows, err := conn.Query(query, args...)
	if err != nil {
		moduleLog.Panicf("Cannot get malls rows: %s", err)
	}
	defer rows.Close()
	var malls []*Mall
	for rows.Next() {
		m := Mall{}
		err = rows.Scan(&m.ID, &m.Name, &m.Phone, &m.LogoSmall, &m.LogoLarge, &m.Location.Lat, &m.Location.Lon, &m.ShopsCount)
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
func GetShopDetails(shopID int, location *Location, cityID *int) *Shop {
	shop := Shop{}
	var err error
	conn := db.GetConnection()
	if location == nil {
		err = conn.QueryRow(`
		SELECT
		  s.id,
		  s.name,
		  s.logo_small,
		  s.logo_large,
		  s.score,
		  s.malls_count,
		  s.phone,
		  s.site
		FROM shop s
		WHERE s.id = $1
		`, shopID).Scan(&shop.ID, &shop.Name, &shop.LogoSmall, &shop.LogoLarge, &shop.Score, &shop.MallsCount, &shop.Phone, &shop.Site)
	} else {
		err = conn.QueryRow(`
		SELECT
		  s.id,
		  s.name,
		  s.logo_small,
		  s.logo_large,
		  s.score,
		  s.malls_count,
		  s.phone,
		  s.site,
		  m.id nearest_mall
		FROM shop s
		  JOIN mall_shop ms ON s.id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.id
		WHERE s.id = $1
		ORDER BY m.location <-> ST_SetSRID(ST_Point($2, $3), 4326)
		LIMIT 1
		`, shopID, location.Lat, location.Lon).Scan(&shop.ID, &shop.Name, &shop.LogoSmall, &shop.LogoLarge, &shop.Score, &shop.MallsCount,
			&shop.Phone, &shop.Site, &shop.NearestMall)
	}
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("shop", shopID).Panicf("Cannot get shop by ID: %s", err)
	}
	return &shop
}
func GetShops(cityID *int, sortKey *string, limit, offset *uint) ([]*Shop, int) {
	var shops []*Shop
	var totalCount int
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		shops = shopsQuery(orderBy.Compile(`
		SELECT
		  s.id,
		  s.name,
		  s.logo_small,
		  s.logo_large,
		  s.score,
		  s.malls_count
		FROM shop s
		  JOIN mall_shop ms ON s.id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.id
		WHERE m.city_id = $3
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`), limit, offset, *cityID)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM shop s
		  JOIN mall_shop ms ON s.id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.id
		WHERE m.city_id = $1
		`, *cityID)
	} else {
		shops = shopsQuery(orderBy.Compile(`
		SELECT
		  s.id,
		  s.name,
		  s.logo_small,
		  s.logo_large,
		  s.score,
		  s.malls_count
		FROM shop s
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`), limit, offset)
		totalCount = countQuery(`
		SELECT
		  count(*) total_count
		FROM shop s
		`)
	}
	return shops, totalCount
}
func GetShopsByMall(mallID int, sortKey *string, limit, offset *uint) ([]*Shop, int) {
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	shops := shopsQuery(orderBy.Compile(`
	SELECT
	  s.id,
	  s.name,
	  s.logo_small,
	  s.logo_large,
	  s.score,
	  s.malls_count
	FROM shop s
	  JOIN mall_shop ms ON s.id = ms.shop_id
	WHERE ms.mall_id = $3
	ORDER BY %s
	LIMIT $1
	OFFSET $2
	`), limit, offset, mallID)
	totalCount := countQuery(`
	SELECT
	  count(*) total_count
	FROM shop s
	  JOIN mall_shop ms ON s.id = ms.shop_id
	WHERE ms.mall_id = $1
	`, mallID)
	return shops, totalCount
}
func GetShopsByIds(shopIDs []int, cityID *int) ([]*Shop, int) {
	if len(shopIDs) == 0 {
		return nil, 0
	}
	shopIDsArray := pq.Array(shopIDs)
	shops := shopsQuery(`
	SELECT
	  s.id,
	  s.name,
	  s.logo_small,
	  s.logo_large,
	  s.score,
	  s.malls_count
	FROM shop s
	WHERE s.id = ANY($1)
	`, shopIDsArray)
	totalCount := countQuery(`
	SELECT
	  count(*) total_count
	FROM shop s
	WHERE s.id = ANY($1)
	`, shopIDsArray)
	return shops, totalCount
}
func GetShopsByName(name string, cityID *int, sortKey *string, limit, offset *uint) ([]*Shop, int) {
	var shops []*Shop
	var totalCount int
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		shops = shopsQuery(orderBy.Compile(`
		SELECT *
		FROM (SELECT DISTINCT ON (s.id)
				s.id,
				s.name,
				s.logo_small,
				s.logo_large,
				s.score,
				s.malls_count
			  FROM shop s
				JOIN shop_name sn ON s.id = sn.shop_id
				JOIN mall_shop ms ON s.id = ms.shop_id
				JOIN mall m ON ms.mall_id = m.id
			  WHERE sn.name ILIKE '%%' || $3 || '%%' AND m.city_id = $4) s
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`), limit, offset, name, cityID)
		totalCount = countQuery(`
		SELECT count(DISTINCT s.id) AS total_count
		FROM shop s
		  JOIN shop_name sn ON s.id = sn.shop_id
		  JOIN mall_shop ms ON s.id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.id
		WHERE sn.name ILIKE '%' || $1 || '%' AND m.city_id = $2
		`, name, cityID)
	} else {
		shops = shopsQuery(orderBy.Compile(`
		SELECT *
		FROM (SELECT DISTINCT ON (s.id)
				s.id,
				s.name,
				s.logo_small,
				s.logo_large,
				s.score,
				s.malls_count
			  FROM shop s
				JOIN shop_name sn ON s.id = sn.shop_id
			  WHERE sn.name ILIKE '%%' || $3 || '%%') s
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`), limit, offset, name)
		totalCount = countQuery(`
		SELECT count(DISTINCT s.id) total_count
		FROM shop s
		  JOIN shop_name sn ON s.id = sn.shop_id
		WHERE sn.name ILIKE '%' || $1 || '%'
		`, name)
	}
	return shops, totalCount
}
func GetShopsByCategory(categoryID int, cityID *int, sortKey *string, limit, offset *uint) ([]*Shop, int) {
	var shops []*Shop
	var totalCount int
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	if cityID != nil {
		shops = shopsQuery(orderBy.Compile(`
		SELECT *
		FROM (SELECT DISTINCT ON (s.id)
				s.id,
				s.name,
				s.logo_small,
				s.logo_large,
				s.score,
				s.malls_count
			  FROM shop s
				JOIN shop_category sc ON s.id = sc.shop_id
				JOIN mall_shop ms ON s.id = ms.shop_id
				JOIN mall m ON ms.mall_id = m.id
			  WHERE sc.category_id = $3 AND m.city_id = $4) s
		ORDER BY %s
		LIMIT $1
		OFFSET $2
	`), limit, offset, categoryID, *cityID)
		totalCount = countQuery(`
		SELECT count(DISTINCT s.id) total_count
		FROM shop s
		  JOIN shop_category sc ON s.id = sc.shop_id
		  JOIN mall_shop ms ON s.id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.id
		WHERE sc.category_id = $1 AND m.city_id = $2
		`, categoryID, *cityID)
	} else {
		shops = shopsQuery(orderBy.Compile(`
		SELECT
		  s.id,
		  s.name,
		  s.logo_small,
		  s.logo_large,
		  s.score,
		  s.malls_count
		FROM shop s
		  JOIN shop_category sc ON s.id = sc.shop_id
		WHERE sc.category_id = $3
		ORDER BY %s
		LIMIT $1
		OFFSET $2
		`), limit, offset, categoryID)
		totalCount = countQuery(`
		SELECT count(*) total_count
		FROM shop s
		  JOIN shop_category sc ON s.id = sc.shop_id
		WHERE sc.category_id = $1
		`, categoryID)
	}
	return shops, totalCount
}
func shopsQuery(query string, args ...interface{}) []*Shop {
	conn := db.GetConnection()
	rows, err := conn.Query(query, args...)
	if err != nil {
		moduleLog.Panicf("Cannot get shops rows: %s", err)
	}
	defer rows.Close()
	var shops []*Shop
	for rows.Next() {
		s := Shop{}
		err = rows.Scan(&s.ID, &s.Name, &s.LogoSmall, &s.LogoLarge, &s.Score, &s.MallsCount)
		if err != nil {
			moduleLog.Panicf("Error during scaning shop row: %s", err)
		}
		shops = append(shops, &s)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning shops rows: %s", err)
	}
	return shops
}
func GetCategoryDetails(categoryID int, cityID *int) *Category {
	conn := db.GetConnection()
	category := Category{}
	err := conn.QueryRow(`
	SELECT
	  c.id,
	  c.name,
	  c.logo_small,
	  c.logo_large,
	  c.shops_count
	FROM category c
	WHERE c.id = $1
	`, categoryID).Scan(&category.ID, &category.Name, &category.LogoSmall, &category.LogoLarge, &category.ShopsCount)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("category", categoryID).Panicf("Cannot get category by ID: %s", err)
	}
	return &category
}
func GetCategories(cityID *int, sortKey *string) ([]*Category, int) {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	categories := categoriesQuery(orderBy.Compile(`
	SELECT
	  c.id,
	  c.name,
	  c.logo_small,
	  c.logo_large,
	  c.shops_count
	FROM category c
	ORDER BY %s
	`))
	totalCount := countQuery(`
	SELECT count(*) total_count
	FROM category c
	`)
	return categories, totalCount
}
func GetCategoriesByIds(categoryIDs []int, cityID *int) ([]*Category, int) {
	categoryIDsArray := pq.Array(categoryIDs)
	categories := categoriesQuery(`
	SELECT
	  c.id,
	  c.name,
	  c.logo_small,
	  c.logo_large,
	  c.shops_count
	FROM category c
	WHERE c.id = ANY ($1)
	`, categoryIDsArray)
	totalCount := countQuery(`
	SELECT count(*) total_count
	FROM category c
	WHERE c.id = ANY ($1)
	`, categoryIDsArray)
	return categories, totalCount
}
func GetCategoriesByShop(shopID int, cityID *int, sortKey *string) ([]*Category, int) {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	categories := categoriesQuery(orderBy.Compile(`
	SELECT
	  c.id,
	  c.name,
	  c.logo_small,
	  c.logo_large,
	  c.shops_count
	FROM category c
	  JOIN shop_category sc ON c.id = sc.category_id
	WHERE sc.shop_id = $1
	ORDER BY %s
	`), shopID)
	totalCount := countQuery(`
	SELECT count(*) total_count
	FROM category c
	  JOIN shop_category sc ON c.id = sc.category_id
	WHERE sc.shop_id = $1
	`, shopID)
	return categories, totalCount
}
func categoriesQuery(query string, args ...interface{}) []*Category {
	conn := db.GetConnection()
	rows, err := conn.Query(query, args...)
	if err != nil {
		moduleLog.Panicf("Cannot get categories rows: %s", err)
	}
	defer rows.Close()
	var shops []*Category
	for rows.Next() {
		c := Category{}
		err = rows.Scan(&c.ID, &c.Name, &c.LogoSmall, &c.LogoLarge, &c.ShopsCount)
		if err != nil {
			moduleLog.Panicf("Error during scaning category row: %s", err)
		}
		shops = append(shops, &c)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning categories rows: %s", err)
	}
	return shops
}
func GetCities(sortKey *string) ([]*City, int) {
	orderBy := CITIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	cities := citiesQuery(orderBy.Compile(`
	SELECT
	  c.id,
	  c.name
	FROM city c
	ORDER BY %s
	`))
	totalCount := countQuery(`
	SELECT count(*) total_count
	FROM city c
	`)
	return cities, totalCount
}
func GetCitiesByName(name string, sortKey *string) ([]*City, int) {
	orderBy := CITIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	cities := citiesQuery(orderBy.Compile(`
	SELECT
	  c.id,
	  c.name
	FROM city c
	WHERE c.name ILIKE '%%' || $1 || '%%'
	ORDER BY %s
	`), name)
	totalCount := countQuery(`
	SELECT count(*) total_count
	FROM city c
	WHERE c.name ILIKE '%%' || $1 || '%%'
	`, name)
	return cities, totalCount
}
func citiesQuery(query string, args ...interface{}) []*City {
	conn := db.GetConnection()
	rows, err := conn.Query(query, args...)
	if err != nil {
		moduleLog.Panicf("Cannot get cities rows: %s", err)
	}
	defer rows.Close()
	var cities []*City
	for rows.Next() {
		c := City{}
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			moduleLog.Panicf("Error during scaning city row: %s", err)
		}
		cities = append(cities, &c)
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panicf("Error after scaning cities rows: %s", err)
	}
	return cities
}
