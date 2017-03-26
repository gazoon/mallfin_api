package models

import (
	"github.com/pkg/errors"
	"mallfin_api/db"
	"mallfin_api/utils"
)

func IsShopExists(shopID int) (bool, error) {
	queryName := utils.CurrentFuncName()
	exists, err := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM shop
		WHERE shop_id = ?0)
	`, shopID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func IsMallExists(mallID int) (bool, error) {
	queryName := utils.CurrentFuncName()
	exists, err := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM mall
		WHERE mall_id = ?0)
	`, mallID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func IsCityExists(cityID int) (bool, error) {
	queryName := utils.CurrentFuncName()
	exists, err := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM city
		WHERE city_id = ?0)
	`, cityID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func IsCategoryExists(categoryID int) (bool, error) {
	queryName := utils.CurrentFuncName()
	exists, err := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM category
		WHERE category_id = ?0)
	`, categoryID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func IsSubwayStationExists(subwayStationID int) (bool, error) {
	queryName := utils.CurrentFuncName()
	exists, err := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM subway_station
		WHERE station_id = ?0)
	`, subwayStationID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func existsQuery(queryName, query string, args ...interface{}) (bool, error) {
	result := struct{ Exists bool }{}
	client := db.GetClient()
	_, err := client.QueryOne(&result, query, args...)
	if err != nil {
		return false, errors.WithMessage(err, queryName)
	}
	return result.Exists, nil
}
