package db

import (
	"mallfin_api/models"
	"mallfin_api/utils"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
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

func (mr *mallRow) toModel() *models.Mall {
	var station *models.SubwayStation
	if mr.StationID != nil && mr.StationName != nil {
		station = &models.SubwayStation{ID: *mr.StationID, Name: *mr.StationName}
	}
	mall := &models.Mall{
		ID:          mr.MallID,
		Name:        mr.MallName,
		Phone:       mr.MallPhone,
		Logo:        models.Logo{Small: mr.MallLogoSmall, Large: mr.MallLogoLarge},
		Location:    models.Location{Lon: mr.MallLocationLon, Lat: mr.MallLocationLat},
		ShopsCount:  mr.ShopsCount,
		Address:     mr.Address,
		DayAndNight: mr.DayAndNight,
		Site:        mr.MallSite,
		Subway:      station,
	}
	return mall
}

func GetMallDetails(mallID int) (*models.Mall, error) {
	queryName := utils.CurrentFuncName()
	mall, err := mallQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
	WHERE m.mall_id = ?0
	LIMIT 1
	`), mallID)
	if err != nil {
		return nil, err
	}
	return mall, nil
}

func GetMallByLocation(location *models.Location) (*models.Mall, error) {
	queryName := utils.CurrentFuncName()
	mall, err := mallQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
	WHERE st_dwithin(st_transform(m.mall_location, 26986), st_transform(ST_Setsrid(st_point(?0, ?1), 4326), 26986), m.mall_radius)
	ORDER BY m.mall_location <-> ST_SetSRID(ST_Point(?0, ?1), 4326)
	LIMIT 1
	`), location.Lon, location.Lat)
	if err != nil {
		return nil, err
	}
	return mall, nil
}

func GetMalls(cityID int, sorting models.Sorting, limit, offset *int) ([]*models.Mall, error) {
	orderBy := mallOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	WHERE m.city_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, cityID)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetMallsWithoutCity(sorting models.Sorting, limit, offset *int) ([]*models.Mall, error) {
	orderBy := mallOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetMallsByIDs(mallIDs []int) ([]*models.Mall, error) {
	if len(mallIDs) == 0 {
		return nil, nil
	}
	mallIDsArray := pg.Array(mallIDs)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM mall m
	WHERE m.mall_id = ANY(?0)
	`), mallIDsArray)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetMallsBySubwayStation(subwayStationID int, sorting models.Sorting, limit, offset *int) ([]*models.Mall, error) {
	orderBy := mallOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	  LEFT JOIN subway_station ss ON m.subway_station_id = ss.station_id
	WHERE ss.station_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, subwayStationID)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetMallsByShop(shopID, cityID int, sorting models.Sorting, limit, offset *int) ([]*models.Mall, error) {
	orderBy := mallOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	  JOIN mall_shop ms ON m.mall_id = ms.mall_id
	WHERE ms.shop_id = ?2 AND m.city_id = ?3
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, shopID, cityID)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetMallsByShopWithoutCity(shopID int, sorting models.Sorting, limit, offset *int) ([]*models.Mall, error) {
	orderBy := mallOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	  JOIN mall_shop ms ON m.mall_id = ms.mall_id
	WHERE ms.shop_id = ?2
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, shopID)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetMallsByName(name string, cityID int, sorting models.Sorting, limit, offset *int) ([]*models.Mall, error) {
	orderBy := mallOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	  JOIN (SELECT DISTINCT ON (mall_id) mall_id
			FROM mall_name
			WHERE mall_name ILIKE '%%' || ?2 || '%%') mn ON m.mall_id = mn.mall_id
	WHERE m.city_id = ?3
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, name, cityID)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetMallsByNameWithoutCity(name string, sorting models.Sorting, limit, offset *int) ([]*models.Mall, error) {
	orderBy := mallOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	malls, err := mallsQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM mall m
	  JOIN (SELECT DISTINCT ON (mall_id) mall_id
			FROM mall_name
			WHERE mall_name ILIKE '%%' || ?2 || '%%') mn ON m.mall_id = mn.mall_id
	ORDER BY {order}
	LIMIT ?0
	OFFSET ?1
	`), limit, offset, name)
	if err != nil {
		return nil, err
	}
	return malls, nil
}

func GetShopsInMalls(mallIDs, shopIDs []int) ([]*models.MallMatchedShops, error) {
	queryName := utils.CurrentFuncName()
	client := GetClient()
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
		return nil, errors.WithMessage(err, queryName)
	}
	mallToShops := map[int][]int{}
	for _, row := range rows {
		mallToShops[row.MallID] = row.Shops
	}
	matchedShops := make([]*models.MallMatchedShops, len(rows))
	for i, row := range rows {
		matchedShops[i] = &models.MallMatchedShops{MallID: row.MallID, ShopIDs: row.Shops}
	}
	for _, mallID := range mallIDs {
		if _, ok := mallToShops[mallID]; !ok {
			matchedShops = append(matchedShops, &models.MallMatchedShops{MallID: mallID, ShopIDs: []int{}})
		}
	}
	return matchedShops, nil
}

func getMallWorkingHours(mallID int) ([]*models.WorkPeriod, error) {
	queryName := utils.CurrentFuncName()
	client := GetClient()
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
		return nil, errors.WithMessage(err, queryName)
	}
	workingHours := make([]*models.WorkPeriod, len(rows))
	for i, row := range rows {
		workingHours[i] = &models.WorkPeriod{
			Open:  models.WeekTime{Day: row.OpenDay, Time: row.OpenTime},
			Close: models.WeekTime{Day: row.CloseDay, Time: row.CloseTime},
		}
	}
	return workingHours, nil
}

func mallQuery(queryName string, queryBasis baseQuery, args ...interface{}) (*models.Mall, error) {
	client := GetClient()
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
		return nil, nil
	} else if err != nil {
		return nil, errors.WithMessage(err, queryName)
	}
	mall := row.toModel()
	if !row.DayAndNight {
		mall.WorkingHours, err = getMallWorkingHours(mall.ID)
		if err != nil {
			return nil, err
		}
	}
	return mall, nil
}

func mallsQuery(queryName string, queryBasis baseQuery, args ...interface{}) ([]*models.Mall, error) {
	client := GetClient()
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
		return nil, errors.WithMessage(err, queryName)
	}
	malls := make([]*models.Mall, len(rows))
	for i, row := range rows {
		malls[i] = row.toModel()
	}
	return malls, nil
}
