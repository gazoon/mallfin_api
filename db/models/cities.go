package models

import (
	"mallfin_api/db"
	"mallfin_api/utils"
)

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
