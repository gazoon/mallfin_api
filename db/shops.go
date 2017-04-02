package db

import (
	"mallfin_api/models"
	"mallfin_api/utils"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
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

func (sr *shopRow) toModel() *models.Shop {
	shop := &models.Shop{
		ID:         sr.ShopID,
		Name:       sr.ShopName,
		Logo:       models.Logo{Small: sr.ShopLogoSmall, Large: sr.ShopLogoLarge},
		Score:      sr.Score,
		MallsCount: sr.MallsCount,
		Phone:      sr.ShopPhone,
		Site:       sr.ShopSite,
	}
	return shop
}

func GetShopDetails(shopID int) (*models.Shop, error) {
	client := GetClient()
	queryName := utils.CurrentFuncName()
	var row shopRow
	_, err := client.QueryOne(&row, `
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
	shop := row.toModel()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.WithMessage(err, queryName)
	}
	return shop, nil
}

func GetShopDetailsWithLocation(shopID int, location *models.Location) (*models.Shop, error) {
	client := GetClient()
	queryName := utils.CurrentFuncName()
	var row struct {
		shopRow
		mallRow
	}
	_, err := client.QueryOne(&row, `
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
	shop := row.shopRow.toModel()
	shop.NearestMall = row.mallRow.toModel()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.WithMessage(err, queryName)
	}
	return shop, nil
}

func GetShops(cityID int, sorting models.Sorting, limit, offset *int) ([]*models.Shop, error) {
	orderBy := shopOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM shop s
	  JOIN mall_shop ms ON s.shop_id = ms.shop_id
	  JOIN mall m ON ms.mall_id = m.mall_id
	WHERE m.city_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, cityID)
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func GetShopsWithoutCity(sorting models.Sorting, limit, offset *int) ([]*models.Shop, error) {
	orderBy := shopOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM shop s
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset)
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func GetShopsByMall(mallID int, sorting models.Sorting, limit, offset *int) ([]*models.Shop, error) {
	orderBy := shopOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM shop s
	  JOIN mall_shop ms ON s.shop_id = ms.shop_id
	WHERE ms.mall_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, mallID)
	if err != nil {
		return nil, nil
	}
	return shops, nil
}

func GetShopsByIDs(shopIDs []int) ([]*models.Shop, error) {
	if len(shopIDs) == 0 {
		return nil, nil
	}
	shopIDsArray := pg.Array(shopIDs)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, `
	SELECT %s
	FROM shop s
	WHERE s.shop_id = ANY(?0)
	`, shopIDsArray)
	if err != nil {
		return nil, nil
	}
	return shops, nil
}

func GetShopsByName(name string, cityID int, sorting models.Sorting, limit, offset *int) ([]*models.Shop, error) {
	orderBy := shopOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, orderBy.CompileBaseQuery(`
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
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func GetShopsByNameWithoutCity(name string, sorting models.Sorting, limit, offset *int) ([]*models.Shop, error) {
	orderBy := shopOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT *
	FROM (SELECT DISTINCT ON (s.shop_id) {columns}
		  FROM shop s
			JOIN shop_name sn ON s.shop_id = sn.shop_id
		  WHERE sn.shop_name ILIKE '%%' || ?2 || '%%') s
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, name)
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func GetShopsByCategory(categoryID, cityID int, sorting models.Sorting, limit, offset *int) ([]*models.Shop, error) {
	orderBy := shopOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, orderBy.CompileBaseQuery(`
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
	`), limit, offset, categoryID, cityID)
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func GetShopsByCategoryWithoutCity(categoryID int, sorting models.Sorting, limit, offset *int) ([]*models.Shop, error) {
	orderBy := shopOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	shops, err := shopsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM shop s
	  JOIN shop_category sc ON s.shop_id = sc.shop_id
	WHERE sc.category_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, categoryID)
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func shopsQuery(queryName string, queryBasis baseQuery, args ...interface{}) ([]*models.Shop, error) {
	client := GetClient()
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
		return nil, errors.WithMessage(err, queryName)
	}
	shops := make([]*models.Shop, len(rows))
	for i, row := range rows {
		shops[i] = row.toModel()
	}
	return shops, nil
}
