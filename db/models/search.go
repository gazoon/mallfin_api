package models

import (
	"mallfin_api/utils"

	"github.com/go-pg/pg"
	"mallfin_api/db"
)

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
