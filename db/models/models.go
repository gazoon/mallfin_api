package models

import (
	"mallfin_api/db"

	"fmt"
	"strings"

	"mallfin_api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/go-pg/pg"
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

func (o *OrderBy) CompileQuery(query string) string {
	return strings.Replace(query, "{order}", o.ToSql(), 1)
}

func (o *OrderBy) CompileBaseQuery(query string) baseQuery {
	return baseQuery(o.CompileQuery(query))
}

type SortKeyToOrderBy struct {
	dict           map[string]*OrderBy
	defaultOrderBy *OrderBy
}

var (
	MALL_DEFAULT_ORDER_BY     = &OrderBy{Column: "m.mall_id", Reverse: false}
	SHOP_DEFAULT_ORDER_BY     = &OrderBy{Column: "s.shop_id", Reverse: false}
	CATEGORY_DEFAULT_ORDER_BY = &OrderBy{Column: "c.category_id", Reverse: false}
	CITY_DEFAULT_ORDER_BY     = &OrderBy{Column: "c.city_id", Reverse: false}
	SEARCH_DEFAULT_ORDER_BY   = &OrderBy{Column: "m.mall_id", Reverse: false}

	MALLS_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           MALL_DEFAULT_ORDER_BY,
			"-id":          {Column: "m.mall_id", Reverse: true},
			"name":         {Column: "m.mall_name", Reverse: false},
			"-name":        {Column: "m.mall_name", Reverse: true},
			"shops_count":  {Column: "m.shops_count", Reverse: false},
			"-shops_count": {Column: "m.shops_count", Reverse: true},
		},
		defaultOrderBy: MALL_DEFAULT_ORDER_BY,
	}
	SHOPS_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           SHOP_DEFAULT_ORDER_BY,
			"-id":          {Column: "s.shop_id", Reverse: true},
			"name":         {Column: "s.shop_name", Reverse: false},
			"-name":        {Column: "s.shop_name", Reverse: true},
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
			"-id":          {Column: "c.cateogry_id", Reverse: true},
			"name":         {Column: "c.category_name", Reverse: false},
			"-name":        {Column: "c.category_name", Reverse: true},
			"shops_count":  {Column: "c.shops_count", Reverse: false},
			"-shops_count": {Column: "c.shops_count", Reverse: true},
		},
		defaultOrderBy: CATEGORY_DEFAULT_ORDER_BY,
	}
	CITIES_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":    CITY_DEFAULT_ORDER_BY,
			"-id":   {Column: "c.city_id", Reverse: true},
			"name":  {Column: "c.city_name", Reverse: false},
			"-name": {Column: "c.city_name", Reverse: true},
		},
		defaultOrderBy: CITY_DEFAULT_ORDER_BY,
	}
	SEARCH_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"mall_id":      SEARCH_DEFAULT_ORDER_BY,
			"-mall_id":     {Column: "m.mall_id", Reverse: true},
			"mall_name":    {Column: "m.mall_name", Reverse: false},
			"-mall_name":   {Column: "m.mall_name", Reverse: true},
			"shops_count":  {Column: "m.shops_count", Reverse: false},
			"-shops_count": {Column: "m.shops_count", Reverse: true},
		},
		defaultOrderBy: SEARCH_DEFAULT_ORDER_BY,
	}
	SEARCH_WITH_DISTANCE_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"mall_id":      SEARCH_DEFAULT_ORDER_BY,
			"-mall_id":     {Column: "m.mall_id", Reverse: true},
			"mall_name":    {Column: "m.mall_name", Reverse: false},
			"-mall_name":   {Column: "m.mall_name", Reverse: true},
			"shops_count":  {Column: "m.shops_count", Reverse: false},
			"-shops_count": {Column: "m.shops_count", Reverse: true},
			"distance":     {Column: "distance", Reverse: false},
			"-distance":    {Column: "distance", Reverse: true},
		},
		defaultOrderBy: SEARCH_DEFAULT_ORDER_BY,
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

type baseQuery string

func (bq baseQuery) withColumns(columns string) string {
	return strings.Replace(string(bq), "{columns}", columns, 1)
}

type WeekTime struct {
	Time string `json:"time"`
	Day  int    `json:"day"`
}
type WorkPeriod struct {
	Open  WeekTime `json:"opening"`
	Close WeekTime `json:"closing"`
}
type Location struct {
	Lat float64
	Lon float64
}
type Logo struct {
	Small string
	Large string
}
type SubwayStation struct {
	ID   int
	Name string
}
type Mall struct {
	ID         int
	Name       string
	Phone      string
	Logo       Logo
	Location   Location
	ShopsCount int
	Address    string
	//Details
	Site         string
	DayAndNight  bool
	Subway       *SubwayStation
	WorkingHours []*WorkPeriod
}
type mallRow struct {
	MallID          int
	MallName        string
	MallPhone       string
	MallLogoLarge   string
	MallLogoSmall   string
	MallLocationLon float64
	MallLocationLat float64
	ShopsCount      int
	Address         string
	//Details
	MallSite    string
	DayAndNight bool
	StationID   *int
	StationName *string
}

func (mr *mallRow) toModel() *Mall {
	var station *SubwayStation
	if mr.StationID != nil && mr.StationName != nil {
		station = &SubwayStation{ID: *mr.StationID, Name: *mr.StationName}
	}
	mall := &Mall{
		ID:          mr.MallID,
		Name:        mr.MallName,
		Phone:       mr.MallPhone,
		Logo:        Logo{Small: mr.MallLogoSmall, Large: mr.MallLogoLarge},
		Location:    Location{Lon: mr.MallLocationLon, Lat: mr.MallLocationLat},
		ShopsCount:  mr.ShopsCount,
		Address:     mr.Address,
		DayAndNight: mr.DayAndNight,
		Site:        mr.MallSite,
		Subway:      station,
	}
	return mall
}

type Shop struct {
	ID         int
	Name       string
	Logo       Logo
	Score      int
	MallsCount int
	//Details
	Phone       string
	Site        string
	NearestMall *Mall
}

type shopRow struct {
	ShopID        int
	ShopName      string
	ShopLogoLarge string
	ShopLogoSmall string
	Score         int
	MallsCount    int
	//Details
	ShopPhone string
	ShopSite  string
}

func (sr *shopRow) toModel() *Shop {
	shop := &Shop{
		ID:         sr.ShopID,
		Name:       sr.ShopName,
		Logo:       Logo{Small: sr.ShopLogoSmall, Large: sr.ShopLogoLarge},
		Score:      sr.Score,
		MallsCount: sr.MallsCount,
		Phone:      sr.ShopPhone,
		Site:       sr.ShopSite,
	}
	return shop
}

type Category struct {
	ID         int
	Name       string
	Logo       Logo
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

func existsQuery(queryName, query string, args ...interface{}) bool {
	result := struct{ Exists bool }{}
	client := db.GetClient()
	_, err := client.QueryOne(&result, query, args...)
	if err != nil {
		moduleLog.WithField("query", queryName).Panicf("Cannot check the existence: %s", err)
	}
	return result.Exists
}
func IsShopExists(shopID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM shop
		WHERE shop_id = ?0)
	`, shopID)
	return exists
}
func IsMallExists(mallID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM mall
		WHERE mall_id = ?0)
	`, mallID)
	return exists
}
func IsCityExists(cityID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM city
		WHERE city_id = ?0)
	`, cityID)
	return exists
}
func IsCategoryExists(categoryID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM category
		WHERE category_id = ?0)
	`, categoryID)
	return exists
}
func IsSubwayStationExists(subwayStationID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM subway_station
		WHERE station_id = ?0)
	`, subwayStationID)
	return exists
}

func searchResultsQuery(queryName, query string, args ...interface{}) []*SearchResult {
	client := db.GetClient()
	locLog := moduleLog.WithField("query", queryName)
	var rows []*struct {
		mallRow
		Shops    []int `pg:",array"`
		Distance *float64
	}
	_, err := client.Query(&rows, query, args...)
	if err != nil {
		locLog.Panicf("Cannot get search results rows: %s", err)
	}
	searchResults := make([]*SearchResult, len(rows))
	for i, row := range rows {
		sr := SearchResult{
			Mall:     row.mallRow.toModel(),
			ShopIDs:  row.Shops,
			Distance: row.Distance,
		}
		searchResults[i] = &sr
	}
	return searchResults
}

func totalCountFromResults(resultsLen int, limit, offset *int) (int, bool) {
	if (limit == nil || *limit == 0) && (offset == nil || *offset == 0 || resultsLen != 0) {
		totalCount := resultsLen
		if offset != nil {
			totalCount += *offset
		}
		return totalCount, true
	}
	return 0, false
}

func GetSearchResults(shopIDs []int, cityID *int, sortKey *string, limit, offset *int) ([]*SearchResult, int) {
	if len(shopIDs) == 0 {
		return nil, 0
	}
	var searchResults []*SearchResult
	var totalCount int
	orderBy := SEARCH_SORT_KEYS.CorrespondingOrderBy(sortKey)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		searchResults = searchResultsQuery(queryName, orderBy.CompileQuery(`
		SELECT
		  m.mall_id,
		  m.mall_name,
		  m.mall_phone,
		  m.mall_logo_small,
		  m.mall_logo_large,
		  ST_Y(m.mall_location) mall_location_lat,
		  ST_X(m.mall_location) mall_location_lon,
		  m.shops_count,
		  m.address,
		  array_agg(ms.shop_id) shops,
		  NULL                  distance
		FROM mall m
		  JOIN mall_shop ms ON m.mall_id = ms.mall_id
		WHERE ms.shop_id = ANY (?2) AND m.city_id = ?3
		GROUP BY m.mall_id
		ORDER BY count(ms.shop_id) DESC, {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, shopIDsArray, *cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0) AND m.city_id = ?1
			`, shopIDsArray, *cityID)
		}
	} else {
		searchResults = searchResultsQuery(queryName, orderBy.CompileQuery(`
		SELECT
		  m.mall_id,
		  m.mall_name,
		  m.mall_phone,
		  m.mall_logo_small,
		  m.mall_logo_large,
		  ST_Y(m.mall_location) mall_location_lat,
		  ST_X(m.mall_location) mall_location_lon,
		  m.shops_count,
		  m.address,
		  array_agg(ms.shop_id) shops,
		  NULL                  distance
		FROM mall m
		  JOIN mall_shop ms ON m.mall_id = ms.mall_id
		WHERE ms.shop_id = ANY (?2)
		GROUP BY m.mall_id
		ORDER BY count(ms.shop_id) DESC, {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, shopIDsArray)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0)
			`, shopIDsArray)
		}
	}
	return searchResults, totalCount
}
func GetSearchResultsWithDistance(shopIDs []int, location *Location, cityID *int, sortKey *string, limit, offset *int) ([]*SearchResult, int) {
	if len(shopIDs) == 0 || location == nil {
		return nil, 0
	}
	var searchResults []*SearchResult
	var totalCount int
	orderBy := SEARCH_WITH_DISTANCE_SORT_KEYS.CorrespondingOrderBy(sortKey)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		searchResults = searchResultsQuery(queryName, orderBy.CompileQuery(`
		SELECT
		  m.mall_id,
		  m.mall_name,
		  m.mall_phone,
		  m.mall_logo_small,
		  m.mall_logo_large,
		  ST_Y(m.mall_location) mall_location_lat,
		  ST_X(m.mall_location) mall_location_lon,
		  m.shops_count,
		  m.address,
		  array_agg(ms.shop_id) shops,
		  st_distance(
			  st_transform(m.mall_location, 26986),
			  st_transform(st_setsrid(st_point(?3, ?4), 4326), 26986)
		  )                     distance
		FROM mall m
		  JOIN mall_shop ms ON m.mall_id = ms.mall_id
		WHERE ms.shop_id = ANY(?2) AND m.city_id = ?5
		GROUP BY m.mall_id
		ORDER BY count(ms.shop_id) DESC, {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, shopIDsArray, location.Lon, location.Lat, *cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0) AND m.city_id = ?1
			`, shopIDsArray, *cityID)
		}
	} else {
		searchResults = searchResultsQuery(queryName, orderBy.CompileQuery(`
		SELECT
		  m.mall_id,
		  m.mall_name,
		  m.mall_phone,
		  m.mall_logo_small,
		  m.mall_logo_large,
		  ST_Y(m.mall_location) mall_location_lat,
		  ST_X(m.mall_location) mall_location_lon,
		  m.shops_count,
		  m.address,
		  array_agg(ms.shop_id) shops,
		  st_distance(
			  st_transform(m.mall_location, 26986),
			  st_transform(st_setsrid(st_point(?3, ?4), 4326), 26986)
		  )                     distance
		FROM mall m
		  JOIN mall_shop ms ON m.mall_id = ms.mall_id
		WHERE ms.shop_id = ANY (?2)
		GROUP BY m.mall_id
		ORDER BY count(ms.shop_id) DESC, {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, shopIDsArray, location.Lon, location.Lat)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0)
			`, shopIDsArray)
		}
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
	queryName := utils.CurrentFuncName()
	locLog := moduleLog.WithField("query", queryName)
	client := db.GetClient()
	var rows []*struct {
		MallID int
		ShopID int
	}
	_, err := client.Query(&rows, `
	SELECT
	  mall_id,
	  shop_id
	FROM mall_shop
	WHERE mall_id = ANY (?0) AND shop_id = ANY (?1)
	`, pg.Array(mallIDs), pg.Array(shopIDs))
	if err != nil && err != pg.ErrNoRows {
		locLog.Panicf("Cannot get shops in malls occurrence: %s", err)
	}
	for _, row := range rows {
		mallsShops[row.MallID] = append(mallsShops[row.MallID], row.ShopID)
	}
	return mallsShops
}
func GetMallWorkingHours(mallID int) []*WorkPeriod {
	queryName := utils.CurrentFuncName()
	locLog := moduleLog.WithFields(log.Fields{"mall": mallID, "query": queryName})
	client := db.GetClient()
	var rows []*struct {
		OpenDay   int
		OpenTime  string
		CloseDay  int
		CloseTime string
	}
	_, err := client.Query(&rows, `
	SELECT
	  open_day,
	  open_time,
	  close_day,
	  close_time
	FROM mall_working_hours
	WHERE mall_id = ?0
	`, mallID)
	if err != nil && err != pg.ErrNoRows {
		locLog.Panicf("Cannot get mall working hours: %s", err)
	}
	workingHours := make([]*WorkPeriod, len(rows))
	for i, row := range rows {
		workingHours[i] = &WorkPeriod{
			Open:  WeekTime{Day: row.OpenDay, Time: row.OpenTime},
			Close: WeekTime{Day: row.CloseDay, Time: row.CloseTime},
		}
	}
	return workingHours
}
func mallQuery(queryName string, queryBasis baseQuery, args ...interface{}) *Mall {
	client := db.GetClient()
	var row mallRow
	query := queryBasis.withColumns(`
	  m.mall_id,
	  m.mall_name,
	  m.mall_phone,
	  m.mall_logo_small,
	  m.mall_logo_large,
	  ST_Y(m.mall_location) mall_location_lat,
	  ST_X(m.mall_location) mall_location_lon,
	  m.shops_count,
	  m.address,
	  m.mall_site,
	  m.day_and_night,
	  ss.station_id,
	  ss.station_name
	`)
	_, err := client.QueryOne(&row, query, args...)
	if err == pg.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithField("query", queryName).Panicf("Cannot get mall: %s", err)
	}
	mall := row.toModel()
	if !row.DayAndNight {
		mall.WorkingHours = GetMallWorkingHours(mall.ID)
	}
	return mall
}
func GetMallDetails(mallID int) *Mall {
	queryName := utils.CurrentFuncName()
	mall := mallQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
	WHERE m.mall_id = ?0
	LIMIT 1
	`), mallID)
	return mall
}
func GetMallByLocation(location *Location) *Mall {
	if location == nil {
		return nil
	}
	queryName := utils.CurrentFuncName()
	mall := mallQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
	WHERE st_dwithin(st_transform(m.mall_location, 26986), st_transform(ST_Setsrid(st_point(?0, ?1), 4326), 26986), m.mall_radius)
	ORDER BY m.mall_location <-> ST_SetSRID(ST_Point(?0, ?1), 4326)
	LIMIT 1
	`), location.Lon, location.Lat)
	return mall
}
func countQuery(queryName, query string, args ...interface{}) int {
	var row struct {
		Count int
	}
	client := db.GetClient()
	_, err := client.QueryOne(&row, query, args...)
	if err != nil {
		moduleLog.WithField("query", queryName).Panicf("Cannot do count query: %s", err)
	}
	return row.Count
}
func GetMalls(cityID *int, sortKey *string, limit, offset *int) ([]*Mall, int) {
	var malls []*Mall
	var totalCount int
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		malls = mallsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM mall m
		WHERE m.city_id = ?2
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, *cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(malls), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM mall m
			WHERE m.city_id = ?0
			`, *cityID)
		}
	} else {
		malls = mallsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM mall m
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(malls), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM mall m
			`)
		}
	}
	return malls, totalCount
}
func GetMallsByIDs(mallIDs []int) ([]*Mall, int) {
	if len(mallIDs) == 0 {
		return nil, 0
	}
	mallIDsArray := pg.Array(mallIDs)
	queryName := utils.CurrentFuncName()
	malls := mallsQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM mall m
	WHERE m.mall_id = ANY(?0)
	`), mallIDsArray)
	totalCount := len(malls)
	return malls, totalCount
}
func GetMallsBySubwayStation(subwayStationID int, sortKey *string, limit, offset *int) ([]*Mall, int) {
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	malls := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
	WHERE ss.station_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, subwayStationID)
	totalCount, ok := totalCountFromResults(len(malls), limit, offset)
	if !ok {
		totalCount = countQuery(queryName, `
		SELECT count(*)
		FROM mall m
		  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
		WHERE ss.station_id = ?0
		`, subwayStationID)
	}
	return malls, totalCount
}
func GetMallsByShop(shopID int, cityID *int, sortKey *string, limit, offset *int) ([]*Mall, int) {
	var malls []*Mall
	var totalCount int
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		malls = mallsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM mall m
		  JOIN mall_shop ms ON m.mall_id = ms.mall_id
		WHERE ms.shop_id = ?2 AND m.city_id = ?3
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, shopID, *cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(malls), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ?0 AND m.city_id = ?1
			`, shopID, *cityID)
		}
	} else {
		malls = mallsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM mall m
		  JOIN mall_shop ms ON m.mall_id = ms.mall_id
		WHERE ms.shop_id = ?2
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, shopID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(malls), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ?0
			`, shopID)
		}
	}
	return malls, totalCount
}
func GetMallsByName(name string, cityID *int, sortKey *string, limit, offset *int) ([]*Mall, int) {
	var malls []*Mall
	var totalCount int
	orderBy := MALLS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		malls = mallsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM mall m
		  JOIN (SELECT DISTINCT ON (mall_id) mall_id
				FROM mall_name
				WHERE mall_name ILIKE '%%' || ?2 || '%%') mn ON m.mall_id = mn.mall_id
		WHERE m.city_id = ?3
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, name, *cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(malls), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM mall m
			  JOIN (SELECT DISTINCT ON (mall_id) mall_id
					FROM mall_name
					WHERE mall_name ILIKE '%' || ?0 || '%') mn ON m.mall_id = mn.mall_id
			WHERE m.city_id = ?1
			`, name, *cityID)
		}
	} else {
		malls = mallsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM mall m
		  JOIN (SELECT DISTINCT ON (mall_id) mall_id
				FROM mall_name
				WHERE mall_name ILIKE '%%' || ?2 || '%%') mn ON m.mall_id = mn.mall_id
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, name)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(malls), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM mall m
			  JOIN (SELECT DISTINCT ON (mall_id) mall_id
					FROM mall_name
					WHERE mall_name ILIKE '%' || ?0 || '%') mn ON m.mall_id = mn.mall_id
			`, name)
		}
	}
	return malls, totalCount
}
func mallsQuery(queryName string, queryBasis baseQuery, args ...interface{}) []*Mall {
	client := db.GetClient()
	locLog := moduleLog.WithField("query", queryName)
	var rows []*mallRow
	query := queryBasis.withColumns(`
	  m.mall_id,
	  m.mall_name,
	  m.mall_phone,
	  m.mall_logo_small,
	  m.mall_logo_large,
	  ST_Y(m.mall_location) mall_location_lat,
	  ST_X(m.mall_location) mall_location_lon,
	  m.shops_count
	`)
	_, err := client.Query(&rows, query, args...)
	if err != nil {
		locLog.Panicf("Cannot get malls rows: %s", err)
	}
	malls := make([]*Mall, len(rows))
	for i, row := range rows {
		malls[i] = row.toModel()
	}
	return malls
}
func GetShopDetails(shopID int, location *Location, cityID *int) *Shop {
	client := db.GetClient()
	queryName := utils.CurrentFuncName()
	var err error
	var shop *Shop
	if location == nil {
		var row shopRow
		_, err = client.QueryOne(&row, `
		SELECT
		  s.shop_id,
		  s.shop_name,
		  s.shop_logo_small,
		  s.shop_logo_large,
		  s.score,
		  s.malls_count,
		  s.shop_phone,
		  s.shop_site
		FROM shop s
		WHERE s.shop_id = ?0
		LIMIT 1
		`, shopID)
		shop = row.toModel()
	} else {
		var row struct {
			shopRow
			mallRow
		}
		_, err = client.QueryOne(&row, `
		SELECT
		  s.shop_id,
		  s.shop_name,
		  s.shop_logo_small,
		  s.shop_logo_large,
		  s.score,
		  s.malls_count,
		  s.shop_phone,
		  s.shop_site,
		  m.mall_id,
		  m.mall_name,
		  m.mall_phone,
		  m.mall_logo_small,
		  m.mall_logo_large,
		  ST_Y(m.mall_location) mall_location_lat,
		  ST_X(m.mall_location) mall_location_lon,
		  m.shops_count
		FROM shop s
		  JOIN mall_shop ms ON s.shop_id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.mall_id
		WHERE s.shop_id = ?0
		ORDER BY m.mall_location <-> ST_SetSRID(ST_Point(?1, ?2), 4326)
		LIMIT 1
		`, shopID, location.Lon, location.Lat)
		shop = row.shopRow.toModel()
		shop.NearestMall = row.mallRow.toModel()
	}
	if err == pg.ErrNoRows {
		return nil
	} else if err != nil {
		moduleLog.WithFields(log.Fields{"shop": shopID, "query": queryName}).Panicf("Cannot get shop by ID: %s", err)
	}
	return shop
}
func GetShops(cityID *int, sortKey *string, limit, offset *int) ([]*Shop, int) {
	var shops []*Shop
	var totalCount int
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		shops = shopsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM shop s
		  JOIN mall_shop ms ON s.shop_id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.mall_id
		WHERE m.city_id = ?2
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, *cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(shops), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM shop s
			  JOIN mall_shop ms ON s.shop_id = ms.shop_id
			  JOIN mall m ON ms.mall_id = m.mall_id
			WHERE m.city_id = ?0
			`, *cityID)
		}
	} else {
		shops = shopsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM shop s
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(shops), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM shop s
			`)
		}
	}
	return shops, totalCount
}
func GetShopsByMall(mallID int, sortKey *string, limit, offset *int) ([]*Shop, int) {
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	shops := shopsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM shop s
	  JOIN mall_shop ms ON s.shop_id = ms.shop_id
	WHERE ms.mall_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, mallID)
	totalCount, ok := totalCountFromResults(len(shops), limit, offset)
	if !ok {
		totalCount = countQuery(queryName, `
		SELECT count(*)
		FROM shop s
		  JOIN mall_shop ms ON s.shop_id = ms.shop_id
		WHERE ms.mall_id = ?0
		`, mallID)
	}
	return shops, totalCount
}
func GetShopsByIDs(shopIDs []int, cityID *int) ([]*Shop, int) {
	if len(shopIDs) == 0 {
		return nil, 0
	}
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	shops := shopsQuery(queryName, `
	SELECT %s
	FROM shop s
	WHERE s.shop_id = ANY(?0)
	`, shopIDsArray)
	totalCount := len(shops)
	return shops, totalCount
}
func GetShopsByName(name string, cityID *int, sortKey *string, limit, offset *int) ([]*Shop, int) {
	var shops []*Shop
	var totalCount int
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		shops = shopsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT *
		FROM (SELECT DISTINCT ON (s.shop_id) {columns}
			  FROM shop s
				JOIN shop_name sn ON s.shop_id = sn.shop_id
				JOIN mall_shop ms ON s.shop_id = ms.shop_id
				JOIN mall m ON ms.mall_id = m.mall_id
			  WHERE sn.shop_name ILIKE '%%' || ?2 || '%%' AND m.city_id = ?3) s
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, name, cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(shops), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(DISTINCT s.shop_id)
			FROM shop s
			  JOIN shop_name sn ON s.shop_id = sn.shop_id
			  JOIN mall_shop ms ON s.shop_id = ms.shop_id
			  JOIN mall m ON ms.mall_id = m.mall_id
			WHERE sn.shop_name ILIKE '%' || ?0 || '%' AND m.city_id = ?1
			`, name, cityID)
		}
	} else {
		shops = shopsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT *
		FROM (SELECT DISTINCT ON (s.shop_id) {columns}
			  FROM shop s
				JOIN shop_name sn ON s.shop_id = sn.shop_id
			  WHERE sn.shop_name ILIKE '%%' || ?2 || '%%') s
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, name)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(shops), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(DISTINCT s.shop_id)
			FROM shop s
			  JOIN shop_name sn ON s.shop_id = sn.shop_id
			WHERE sn.shop_name ILIKE '%' || ?0 || '%'
			`, name)
		}
	}
	return shops, totalCount
}
func GetShopsByCategory(categoryID int, cityID *int, sortKey *string, limit, offset *int) ([]*Shop, int) {
	var shops []*Shop
	var totalCount int
	orderBy := SHOPS_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		shops = shopsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT *
		FROM (SELECT DISTINCT ON (s.shop_id) {columns}
			  FROM shop s
				JOIN shop_category sc ON s.shop_id = sc.shop_id
				JOIN mall_shop ms ON s.shop_id = ms.shop_id
				JOIN mall m ON ms.mall_id = m.mall_id
			  WHERE sc.category_id = ?2 AND m.city_id = ?3) s
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
	`), limit, offset, categoryID, *cityID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(shops), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(DISTINCT s.shop_id)
			FROM shop s
			  JOIN shop_category sc ON s.shop_id = sc.shop_id
			  JOIN mall_shop ms ON s.shop_id = ms.shop_id
			  JOIN mall m ON ms.mall_id = m.mall_id
			WHERE sc.category_id = ?0 AND m.city_id = ?1
			`, categoryID, *cityID)
		}
	} else {
		shops = shopsQuery(queryName, orderBy.CompileBaseQuery(`
		SELECT {columns}
		FROM shop s
		  JOIN shop_category sc ON s.shop_id = sc.shop_id
		WHERE sc.category_id = ?2
		ORDER BY {order}
		LIMIT ?0
		OFFSET ?1
		`), limit, offset, categoryID)
		var ok bool
		if totalCount, ok = totalCountFromResults(len(shops), limit, offset); !ok {
			totalCount = countQuery(queryName, `
			SELECT count(*)
			FROM shop s
			  JOIN shop_category sc ON s.shop_id = sc.shop_id
			WHERE sc.category_id = ?0
			`, categoryID)
		}
	}
	return shops, totalCount
}
func shopsQuery(queryName string, queryBasis baseQuery, args ...interface{}) []*Shop {
	client := db.GetClient()
	locLog := moduleLog.WithField("query", queryName)
	query := queryBasis.withColumns(`
	  s.shop_id,
	  s.shop_name,
	  s.shop_logo_small,
	  s.shop_logo_large,
	  s.score,
	  s.malls_count
	`)
	var rows []*shopRow
	_, err := client.Query(&rows, query, args...)
	if err != nil {
		locLog.Panicf("Cannot get shops rows: %s", err)
	}
	shops := make([]*Shop, len(rows))
	for i, row := range rows {
		shops[i] = row.toModel()
	}
	return shops
}
func GetCategoryDetails(categoryID int, cityID *int) *Category {
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM category c
	WHERE c.category_id = ?0
	LIMIT 1
	`), categoryID)
	if len(categories) == 0 {
		return nil
	}
	return categories[0]
}
func GetCategories(cityID *int, sortKey *string) []*Category {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM category c
	ORDER BY {order}
	`))
	return categories
}
func GetCategoriesByIDs(categoryIDs []int, cityID *int) ([]*Category, int) {
	categoryIDsArray := pg.Array(categoryIDs)
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM category c
	WHERE c.category_id = ANY (?0)
	`), categoryIDsArray)
	totalCount := len(categories)
	return categories, totalCount
}
func GetCategoriesByShop(shopID int, cityID *int, sortKey *string) []*Category {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM category c
	  JOIN shop_category sc ON c.category_id = sc.category_id
	WHERE sc.shop_id = ?0
	ORDER BY {order}
	`), shopID)
	return categories
}
func categoriesQuery(queryName string, queryBasis baseQuery, args ...interface{}) []*Category {
	client := db.GetClient()
	query := queryBasis.withColumns(`
	  c.category_id,
	  c.category_name,
	  c.category_logo_small,
	  c.category_logo_large,
	  c.shops_count
	`)
	var rows []*struct {
		CategoryID        int
		CategoryName      string
		CategoryLogoSmall string
		CategoryLogoLarge string
		ShopsCount        int
	}
	_, err := client.Query(&rows, query, args...)
	locLog := moduleLog.WithField("query", queryName)
	if err != nil {
		locLog.Panicf("Cannot get categories rows: %s", err)
	}
	categories := make([]*Category, len(rows))
	for i, row := range rows {
		categories[i] = &Category{
			ID:         row.CategoryID,
			Name:       row.CategoryName,
			Logo:       Logo{Small: row.CategoryLogoSmall, Large: row.CategoryLogoLarge},
			ShopsCount: row.ShopsCount,
		}
	}
	return categories
}
func GetCities(sortKey *string) []*City {
	orderBy := CITIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	cities := citiesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM city c
	ORDER BY {order}
	`))
	return cities
}
func GetCitiesByName(name string, sortKey *string) []*City {
	orderBy := CITIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	cities := citiesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM city c
	WHERE c.city_name ILIKE '%%' || ?0 || '%%'
	ORDER BY {order}
	`), name)
	return cities
}
func GetCityByLocation(location *Location) *City {
	if location == nil {
		return nil
	}
	queryName := utils.CurrentFuncName()
	cities := citiesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM city c
	WHERE st_dwithin(st_transform(c.city_location, 26986), st_transform(ST_Setsrid(st_point(?, ?), 4326), 26986), c.city_radius)
	ORDER BY c.city_location <-> ST_SetSRID(ST_Point(?, ?), 4326)
	LIMIT 1
	`), location.Lon, location.Lat, location.Lon, location.Lat)
	if len(cities) == 0 {
		return nil
	}
	return cities[0]
}
func citiesQuery(queryName string, queryBasis baseQuery, args ...interface{}) []*City {
	client := db.GetClient()
	locLog := moduleLog.WithField("query", queryName)
	var rows []*struct {
		CityID   int
		CityName string
	}
	query := queryBasis.withColumns(`
	  c.city_id,
	  c.city_name
	`)
	_, err := client.Query(&rows, query, args...)
	if err != nil {
		locLog.Panicf("Cannot get cities rows: %s", err)
	}
	cities := make([]*City, len(rows))
	for i, row := range rows {
		cities[i] = &City{ID: row.CityID, Name: row.CityName}
	}
	return cities
}
