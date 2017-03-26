package models

import (
	"mallfin_api/utils"
	"mallfin_api/db"
	"github.com/go-pg/pg"
	log "github.com/Sirupsen/logrus"
)

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

func GetShopsInMalls(mallIDs, shopIDs []int) []*MallMatchedShops {
	queryName := utils.CurrentFuncName()
	locLog := moduleLog.WithField("query", queryName)
	client := db.GetClient()
	var rows []*struct {
		MallID int
		Shops  []int `pg:",array"`
	}
	_, err := client.Query(&rows, `
	SELECT
	  mall_id,
	  array_agg(shop_id) shops
	FROM mall_shop
	WHERE mall_id = ANY (?) AND shop_id = ANY (?)
	GROUP BY mall_id
	ORDER BY count(shop_id) DESC
	`, pg.Array(mallIDs), pg.Array(shopIDs))
	if err != nil && err != pg.ErrNoRows {
		locLog.Panicf("Cannot get shops in malls occurrence: %s", err)
	}
	mallToShops := map[int][]int{}
	for _, row := range rows {
		mallToShops[row.MallID] = row.Shops
	}
	matchedShops := make([]*MallMatchedShops, len(rows))
	for i, row := range rows {
		matchedShops[i] = &MallMatchedShops{MallID: row.MallID, ShopIDs: row.Shops}
	}
	for _, mallID := range mallIDs {
		if _, ok := mallToShops[mallID]; !ok {
			matchedShops = append(matchedShops, &MallMatchedShops{MallID: mallID, ShopIDs: []int{}})
		}
	}
	return matchedShops
}

func getMallWorkingHours(mallID int) []*WorkPeriod {
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
		mall.WorkingHours = getMallWorkingHours(mall.ID)
	}
	return mall
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
