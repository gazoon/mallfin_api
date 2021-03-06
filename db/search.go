package db

import (
	"mallfin_api/utils"

	"mallfin_api/models"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
)

func GetSearchResults(shopIDs []int, cityID int, sorting models.Sorting, limit, offset *int) ([]*models.SearchResult, error) {
	if len(shopIDs) == 0 {
		return nil, nil
	}
	orderBy := searchOrderBy(sorting)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	searchResults, err := searchResultsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	  NULL                  distance
	FROM mall m
	  JOIN mall_shop ms ON m.mall_id = ms.mall_id
	WHERE ms.shop_id = ANY (?2) AND m.city_id = ?3
	GROUP BY m.mall_id
	ORDER BY count(ms.shop_id) DESC, {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, shopIDsArray, cityID)
	if err != nil {
		return nil, err
	}
	return searchResults, nil
}

func GetSearchResultsWithoutCity(shopIDs []int, sorting models.Sorting, limit, offset *int) ([]*models.SearchResult, error) {
	if len(shopIDs) == 0 {
		return nil, nil
	}
	orderBy := searchOrderBy(sorting)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	searchResults, err := searchResultsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
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
		return nil, err
	}
	return searchResults, nil
}

func GetSearchResultsWithDistance(shopIDs []int, location *models.Location, cityID int, sorting models.Sorting, limit, offset *int) ([]*models.SearchResult, error) {
	if len(shopIDs) == 0 {
		return nil, nil
	}
	orderBy := searchOrderBy(sorting)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	searchResults, err := searchResultsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
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
	`), limit, offset, shopIDsArray, location.Lon, location.Lat, cityID)
	if err != nil {
		return nil, err
	}
	return searchResults, nil
}

func GetSearchResultsWithDistanceWithoutCity(shopIDs []int, location *models.Location, sorting models.Sorting, limit, offset *int) ([]*models.SearchResult, error) {
	if len(shopIDs) == 0 {
		return nil, nil
	}
	orderBy := searchOrderBy(sorting)
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	searchResults, err := searchResultsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
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
		return nil, err
	}
	return searchResults, nil
}

func searchResultsQuery(queryName string, queryBasis baseQuery, args ...interface{}) ([]*models.SearchResult, error) {
	client := GetClient()
	var rows []*struct {
		mallRow
		Shops    []int `pg:",array"`
		Distance *float64
	}
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
	  array_agg(ms.shop_id) shops,
	`)
	_, err := client.Query(&rows, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, queryName)
	}
	searchResults := make([]*models.SearchResult, len(rows))
	for i, row := range rows {
		sr := models.SearchResult{
			Mall:     row.mallRow.toModel(),
			ShopIDs:  row.Shops,
			Distance: row.Distance,
		}
		searchResults[i] = &sr
	}
	return searchResults, nil
}
