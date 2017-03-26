package models

import (
	"mallfin_api/db"
	"mallfin_api/utils"
)

func IsShopExists(shopID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM shop
		WHERE shop_id = ?0)
	`, shopID)
	return exists
}

func IsMallExists(mallID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM mall
		WHERE mall_id = ?0)
	`, mallID)
	return exists
}

func IsCityExists(cityID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM city
		WHERE city_id = ?0)
	`, cityID)
	return exists
}

func IsCategoryExists(categoryID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM category
		WHERE category_id = ?0)
	`, categoryID)
	return exists
}

func IsSubwayStationExists(subwayStationID int) bool {
	queryName := utils.CurrentFuncName()
	exists := existsQuery(queryName, `
	SELECT exists(
		SELECT *
		FROM subway_station
		WHERE station_id = ?0)
	`, subwayStationID)
	return exists
}

func existsQuery(queryName, query string, args ...interface{}) bool {
	result := struct{ Exists bool }{}
	client := db.GetClient()
	_, err := client.QueryOne(&result, query, args...)
	if err != nil {
		moduleLog.WithField("query", queryName).Panicf("Cannot check the existence: %s", err)
	}
	return result.Exists
}
