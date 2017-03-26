package models

import (
	"mallfin_api/utils"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"mallfin_api/db"
)

func GetSearchResults(shopIDs []int, cityID *int, sortKey *string, limit, offset *int) ([]*SearchResult, int, error) {
	if len(shopIDs) == 0 {
		return nil, 0, nil
	}
	var searchResults []*SearchResult
	var totalCount int
	var err error
	orderBy := SEARCH_SORT_KEYS.CorrespondingOrderBy(sortKey)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		searchResults, err = searchResultsQuery(queryName, orderBy.CompileQuery(`
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
		if err != nil {
			return nil, 0, err
		}
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount, err = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0) AND m.city_id = ?1
			`, shopIDsArray, *cityID)
			if err != nil {
				return nil, 0, err
			}
		}
	} else {
		searchResults, err = searchResultsQuery(queryName, orderBy.CompileQuery(`
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
		if err != nil {
			return nil, 0, err
		}
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount, err = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0)
			`, shopIDsArray)
			if err != nil {
				return nil, 0, err
			}
		}
	}
	return searchResults, totalCount, nil
}

func GetSearchResultsWithDistance(shopIDs []int, location *Location, cityID *int, sortKey *string, limit, offset *int) ([]*SearchResult, int, error) {
	if len(shopIDs) == 0 {
		return nil, 0, nil
	}
	var searchResults []*SearchResult
	var totalCount int
	var err error
	orderBy := SEARCH_WITH_DISTANCE_SORT_KEYS.CorrespondingOrderBy(sortKey)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	if cityID != nil {
		searchResults, err = searchResultsQuery(queryName, orderBy.CompileQuery(`
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
		if err != nil {
			return nil, 0, err
		}
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount, err = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0) AND m.city_id = ?1
			`, shopIDsArray, *cityID)
			if err != nil {
				return nil, 0, err
			}
		}
	} else {
		searchResults, err = searchResultsQuery(queryName, orderBy.CompileQuery(`
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
		if err != nil {
			return nil, 0, err
		}
		var ok bool
		if totalCount, ok = totalCountFromResults(len(searchResults), limit, offset); !ok {
			totalCount, err = countQuery(queryName, `
			SELECT count(DISTINCT m.mall_id)
			FROM mall m
			  JOIN mall_shop ms ON m.mall_id = ms.mall_id
			WHERE ms.shop_id = ANY (?0)
			`, shopIDsArray)
			if err != nil {
				return nil, 0, err
			}
		}
	}
	return searchResults, totalCount, nil
}

func searchResultsQuery(queryName, query string, args ...interface{}) ([]*SearchResult, error) {
	client := db.GetClient()
	var rows []*struct {
		mallRow
		Shops    []int `pg:",array"`
		Distance *float64
	}
	_, err := client.Query(&rows, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, queryName)
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
	return searchResults, nil
}
