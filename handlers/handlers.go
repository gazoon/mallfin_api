package handlers

import (
	"net/http"

	"mallfin_api/db"
	"mallfin_api/models"
	"mallfin_api/serializers"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/binding"
	"github.com/gazoon/httprouter"
)

func MallsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := mallsListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	var malls []*models.Mall
	var totalCount int
	var err error
	if formData.SubwayStation != nil {
		subwayStationID := *formData.SubwayStation
		if !checkSubwayStation(w, subwayStationID) {
			return
		}
		malls, totalCount, err = db.GetMallsBySubwayStation(subwayStationID, sorting, limit, offset)
	} else if formData.Query != nil {
		name := *formData.Query
		malls, totalCount, err = db.GetMallsByName(name, cityID, sorting, limit, offset)
	} else if formData.Shop != nil {
		shopID := *formData.Shop
		if !checkShop(w, shopID) {
			return
		}
		malls, totalCount, err = db.GetMallsByShop(shopID, cityID, sorting, limit, offset)
	} else {
		malls, totalCount, err = db.GetMalls(cityID, sorting, limit, offset)
	}
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(w, r, serialized, totalCount, limit, offset)
}

func MallDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mallID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	mall, err := db.GetMallDetails(mallID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	if mall == nil {
		errorResponse(w, MALL_NOT_FOUND, "Mall with such id does not exists", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeMall(mall)
	response(w, serialized)
}

func ShopsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := shopsListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	var shops []*models.Shop
	var totalCount int
	var err error
	if formData.Mall != nil {
		mallID := *formData.Mall
		if !checkMall(w, mallID) {
			return
		}
		shops, totalCount, err = db.GetShopsByMall(mallID, sorting, limit, offset)
	} else if formData.Query != nil {
		name := *formData.Query
		shops, totalCount, err = db.GetShopsByName(name, cityID, sorting, limit, offset)
	} else if formData.Category != nil {
		categoryID := *formData.Category
		if !checkCategory(w, categoryID) {
			return
		}
		shops, totalCount, err = db.GetShopsByCategory(categoryID, cityID, sorting, limit, offset)
	} else {
		shops, totalCount, err = db.GetShops(cityID, sorting, limit, offset)
	}
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	serialized := serializers.SerializeShops(shops)
	paginateResponse(w, r, serialized, totalCount, limit, offset)
}

func ShopDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	formData := shopDetailsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	shopID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	var userLocation *models.Location = nil
	if formData.LocationLat != nil && formData.LocationLon != nil {
		userLocation = &models.Location{
			Lat: *formData.LocationLat,
			Lon: *formData.LocationLon,
		}
	}
	shop, err := db.GetShopDetails(shopID, userLocation, cityID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	if shop == nil {
		errorResponse(w, SHOP_NOT_FOUND, "Shop with such id does not exists", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeShop(shop)
	response(w, serialized)
}

func CategoriesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := categoriesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	var categories []*models.Category
	var err error
	if formData.Shop != nil {
		shopID := *formData.Shop
		if !checkShop(w, shopID) {
			return
		}
		categories, err = db.GetCategoriesByShop(shopID, cityID, sorting)
	} else {
		categories, err = db.GetCategories(cityID, sorting)
	}
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	serialized := serializers.SerializeCategories(categories)
	response(w, serialized)
}

func CategoryDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	formData := categoryDetailsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	categoryID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	category, err := db.GetCategoryDetails(categoryID, cityID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	if category == nil {
		errorResponse(w, CATEGORY_NOT_FOUND, "Category with such id does not exists", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeCategory(category)
	response(w, serialized)
}

func CitiesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := citiesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	var cities []*models.City
	var err error
	if formData.Query != nil {
		name := *formData.Query
		cities, err = db.GetCitiesByName(name, sorting)
	} else {
		cities, err = db.GetCities(sorting)
	}
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	serialized := serializers.SerializeCities(cities)
	response(w, serialized)
}

func CurrentCity(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := CoordinatesForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	userLocation := &models.Location{
		Lat: formData.LocationLat,
		Lon: formData.LocationLon,
	}
	city, err := db.GetCityByLocation(userLocation)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	if city == nil {
		errorResponse(w, CITY_NOT_FOUND, "In this place there is no city.", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeCity(city)
	response(w, serialized)
}

func CurrentMall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := CoordinatesForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	userLocation := &models.Location{
		Lat: formData.LocationLat,
		Lon: formData.LocationLon,
	}
	mall, err := db.GetMallByLocation(userLocation)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	if mall == nil {
		errorResponse(w, MALL_NOT_FOUND, "In this place there is no mall.", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeMall(mall)
	response(w, serialized)
}

func ShopsInMalls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := shopsInMallsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	mallIDs := formData.Malls
	shopIDs := formData.Shops
	mallsShops, err := db.GetShopsInMalls(mallIDs, shopIDs)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	serialized := serializers.SerializeShopsInMalls(mallsShops)
	response(w, serialized)
}

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := searchForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	shopIDs := formData.Shops
	sorting := formData.Sort
	if !checkCity(w, cityID) {
		return
	}
	var searchResults []*models.SearchResult
	var totalCount int
	var err error
	if formData.LocationLat != nil && formData.LocationLon != nil {
		userLocation := &models.Location{
			Lat: *formData.LocationLat,
			Lon: *formData.LocationLon,
		}
		searchResults, totalCount, err = db.GetSearchResultsWithDistance(shopIDs, userLocation, cityID, sorting, limit, offset)
	} else {
		searchResults, totalCount, err = db.GetSearchResults(shopIDs, cityID, sorting, limit, offset)
	}
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	serialized := serializers.SerializeSearchResults(searchResults)
	paginateResponse(w, r, serialized, totalCount, limit, offset)
}

func checkCity(w http.ResponseWriter, cityID *int) bool {
	if cityID != nil {
		exists, err := db.IsCityExists(*cityID)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return false
		}
		if !exists {
			errorResponse(w, CITY_NOT_FOUND, "City with such id does not exists.", http.StatusNotFound)
			return false
		}
	}
	return true
}

func checkShop(w http.ResponseWriter, shopID int) bool {
	exists, err := db.IsShopExists(shopID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		errorResponse(w, SHOP_NOT_FOUND, "Shop with such id does not exists.", http.StatusNotFound)
		return false
	}
	return true
}

func checkCategory(w http.ResponseWriter, categoryID int) bool {
	exists, err := db.IsCategoryExists(categoryID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		errorResponse(w, CATEGORY_NOT_FOUND, "Category with such id does not exists.", http.StatusNotFound)
		return false
	}
	return true
}

func checkMall(w http.ResponseWriter, mallID int) bool {
	exists, err := db.IsMallExists(mallID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		errorResponse(w, MALL_NOT_FOUND, "Mall with such id does not exists.", http.StatusNotFound)
		return false
	}
	return true
}

func checkSubwayStation(w http.ResponseWriter, subwayStationID int) bool {
	exists, err := db.IsSubwayStationExists(subwayStationID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		errorResponse(w, SUBWAY_STATION_NOT_FOUND, "Subway station with such id does not exists.", http.StatusNotFound)
		return false
	}
	return true
}
