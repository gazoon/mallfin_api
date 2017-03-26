package models

import (
	"github.com/pkg/errors"
	"mallfin_api/db"
	"mallfin_api/utils"
)

func GetCities(sortKey *string) ([]*City, error) {
	orderBy := CITIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	cities, err := citiesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM city c
	ORDER BY {order}
	`))
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func GetCitiesByName(name string, sortKey *string) ([]*City, error) {
	orderBy := CITIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	cities, err := citiesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM city c
	WHERE c.city_name ILIKE '%%' || ?0 || '%%'
	ORDER BY {order}
	`), name)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func GetCityByLocation(location *Location) (*City, error) {
	queryName := utils.CurrentFuncName()
	cities, err := citiesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM city c
	WHERE st_dwithin(st_transform(c.city_location, 26986), st_transform(ST_Setsrid(st_point(?, ?), 4326), 26986), c.city_radius)
	ORDER BY c.city_location <-> ST_SetSRID(ST_Point(?, ?), 4326)
	LIMIT 1
	`), location.Lon, location.Lat, location.Lon, location.Lat)
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, nil
	}
	return cities[0], nil
}

func citiesQuery(queryName string, queryBasis baseQuery, args ...interface{}) ([]*City, error) {
	client := db.GetClient()
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
		return nil, errors.WithMessage(err, queryName)
	}
	cities := make([]*City, len(rows))
	for i, row := range rows {
		cities[i] = &City{ID: row.CityID, Name: row.CityName}
	}
	return cities, nil
}
