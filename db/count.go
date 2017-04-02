package db

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"mallfin_api/utils"
)

func MallsCount(cityID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM mall m
	WHERE m.city_id = ?0
	`, cityID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func MallsWithoutCityCount() (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM mall m
	`)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func MallsByNameCount(name string, cityID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM mall m
	  JOIN (SELECT DISTINCT ON (mall_id) mall_id
			FROM mall_name
			WHERE mall_name ILIKE '%' || ?0 || '%') mn ON m.mall_id = mn.mall_id
	WHERE m.city_id = ?1
	`, name, cityID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func MallsByNameWithoutCityCount(name string) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM mall m
	  JOIN (SELECT DISTINCT ON (mall_id) mall_id
			FROM mall_name
			WHERE mall_name ILIKE '%' || ?0 || '%') mn ON m.mall_id = mn.mall_id
	`, name)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func MallsByShopWithoutCityCount(shopID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM mall m
	  JOIN mall_shop ms ON m.mall_id = ms.mall_id
	WHERE ms.shop_id = ?0
	`, shopID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func MallsByShopCount(shopID, cityID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM mall m
	  JOIN mall_shop ms ON m.mall_id = ms.mall_id
	WHERE ms.shop_id = ?0 AND m.city_id = ?1
	`, shopID, cityID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func MallsBySubwayStationCount(subwayStationID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
	WHERE ss.station_id = ?0
	`, subwayStationID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func SearchResultsWithoutCityCount(shopIDs []int) (int, error) {
	queryName := utils.CurrentFuncName()
	shopIDsArray := pg.Array(shopIDs)
	totalCount, err := countQuery(queryName, `
	SELECT count(DISTINCT m.mall_id)
	FROM mall m
	  JOIN mall_shop ms ON m.mall_id = ms.mall_id
	WHERE ms.shop_id = ANY (?0)
	`, shopIDsArray)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func SearchResultsCount(shopIDs []int, cityID int) (int, error) {
	queryName := utils.CurrentFuncName()
	shopIDsArray := pg.Array(shopIDs)
	totalCount, err := countQuery(queryName, `
	SELECT count(DISTINCT m.mall_id)
	FROM mall m
	  JOIN mall_shop ms ON m.mall_id = ms.mall_id
	WHERE ms.shop_id = ANY (?0) AND m.city_id = ?1
	`, shopIDsArray, cityID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func ShopsByCategoryWithoutCityCount(categoryID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM shop s
	  JOIN shop_category sc ON s.shop_id = sc.shop_id
	WHERE sc.category_id = ?0
	`, categoryID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func ShopsByCategoryCount(categoryID, cityID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(DISTINCT s.shop_id)
	FROM shop s
	  JOIN shop_category sc ON s.shop_id = sc.shop_id
	  JOIN mall_shop ms ON s.shop_id = ms.shop_id
	  JOIN mall m ON ms.mall_id = m.mall_id
	WHERE sc.category_id = ?0 AND m.city_id = ?1
	`, categoryID, cityID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func ShopsByNameWithoutCityCount(name string) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(DISTINCT s.shop_id)
	FROM shop s
	  JOIN shop_name sn ON s.shop_id = sn.shop_id
	WHERE sn.shop_name ILIKE '%' || ?0 || '%'
	`, name)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func ShopsByNameCount(name string, cityID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(DISTINCT s.shop_id)
	FROM shop s
	  JOIN shop_name sn ON s.shop_id = sn.shop_id
	  JOIN mall_shop ms ON s.shop_id = ms.shop_id
	  JOIN mall m ON ms.mall_id = m.mall_id
	WHERE sn.shop_name ILIKE '%' || ?0 || '%' AND m.city_id = ?1
	`, name, cityID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func ShopsByMallCount(mallID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM shop s
	  JOIN mall_shop ms ON s.shop_id = ms.shop_id
	WHERE ms.mall_id = ?0
	`, mallID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func ShopsWithoutCityCount() (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM shop s
	`)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func ShopsCount(cityID int) (int, error) {
	queryName := utils.CurrentFuncName()
	totalCount, err := countQuery(queryName, `
	SELECT count(*)
	FROM shop s
	  JOIN mall_shop ms ON s.shop_id = ms.shop_id
	  JOIN mall m ON ms.mall_id = m.mall_id
	WHERE m.city_id = ?0
	`, cityID)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func countQuery(queryName, query string, args ...interface{}) (int, error) {
	var row struct {
		Count int
	}
	client := GetClient()
	_, err := client.QueryOne(&row, query, args...)
	if err != nil {
		return 0, errors.WithMessage(err, queryName)
	}
	return row.Count, nil
}
