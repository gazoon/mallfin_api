package models

import (
	"mallfin_api/db"
	"mallfin_api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/go-pg/pg"
)

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
