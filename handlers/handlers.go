package handlers

import (
	"net/http"

	"mallfin_api/db/models"
	"mallfin_api/serializers"

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
	sortKey := formData.Sort
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	var malls []*models.Mall
	var totalCount int
	if formData.SubwayStation != nil {
		subwayStationID := *formData.SubwayStation
		if !models.IsSubwayStationExists(subwayStationID) {
			errorResponse(w, SUBWAY_STATION_NOT_FOUND, "Subway station with such id does not exists.", http.StatusNotFound)
			return
		}
		malls, totalCount = models.GetMallsBySubwayStation(subwayStationID, sortKey, limit, offset)
	} else if formData.Query != nil {
		name := *formData.Query
		malls, totalCount = models.GetMallsByName(name, cityID, sortKey, limit, offset)
	} else if formData.Shop != nil {
		shopID := *formData.Shop
		if !models.IsShopExists(shopID) {
			errorResponse(w, SHOP_NOT_FOUND, "Shop with such id does not exists.", http.StatusNotFound)
			return
		}
		malls, totalCount = models.GetMallsByShop(shopID, cityID, sortKey, limit, offset)
	} else {
		malls, totalCount = models.GetMalls(cityID, sortKey, limit, offset)
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
	mall := models.GetMallDetails(mallID)
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
	sortKey := formData.Sort
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	var shops []*models.Shop
	var totalCount int
	if formData.Mall != nil {
		mallID := *formData.Mall
		if !models.IsMallExists(mallID) {
			errorResponse(w, MALL_NOT_FOUND, "Mall with such id does not exists.", http.StatusNotFound)
			return
		}
		shops, totalCount = models.GetShopsByMall(mallID, sortKey, limit, offset)
	} else if formData.Query != nil {
		name := *formData.Query
		shops, totalCount = models.GetShopsByName(name, cityID, sortKey, limit, offset)
	} else if formData.Category != nil {
		categoryID := *formData.Category
		if !models.IsCategoryExists(categoryID) {
			errorResponse(w, CATEGORY_NOT_FOUND, "Category with such id does not exists.", http.StatusNotFound)
			return
		}
		shops, totalCount = models.GetShopsByCategory(categoryID, cityID, sortKey, limit, offset)
	} else {
		shops, totalCount = models.GetShops(cityID, sortKey, limit, offset)
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
	shop := models.GetShopDetails(shopID, userLocation, cityID)
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
	sortKey := formData.Sort
	cityID := formData.City
	if !checkCity(w, cityID) {
		return
	}
	var categories []*models.Category
	if formData.Shop != nil {
		shopID := *formData.Shop
		if !models.IsShopExists(shopID) {
			errorResponse(w, SHOP_NOT_FOUND, "Shop with such id does not exists.", http.StatusNotFound)
			return
		}
		categories = models.GetCategoriesByShop(shopID, cityID, sortKey)
	} else {
		categories = models.GetCategories(cityID, sortKey)
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
	category := models.GetCategoryDetails(categoryID, cityID)
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
	sortKey := formData.Sort
	var cities []*models.City
	if formData.Query != nil {
		name := *formData.Query
		cities = models.GetCitiesByName(name, sortKey)
	} else {
		cities = models.GetCities(sortKey)
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
	city := models.GetCityByLocation(userLocation)
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
	mall := models.GetMallByLocation(userLocation)
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
	mallsShops := models.GetShopsInMalls(mallIDs, shopIDs)
	//serialized := serializers.SerializeShopsInMalls(mallsShops)
	response(w, mallsShops)
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
	sortKey := formData.Sort
	if !checkCity(w, cityID) {
		return
	}
	var searchResults []*models.SearchResult
	var totalCount int
	if formData.LocationLat != nil && formData.LocationLon != nil {
		userLocation := &models.Location{
			Lat: *formData.LocationLat,
			Lon: *formData.LocationLon,
		}
		searchResults, totalCount = models.GetSearchResultsWithDistance(shopIDs, userLocation, cityID, sortKey, limit, offset)
	} else {
		searchResults, totalCount = models.GetSearchResults(shopIDs, cityID, sortKey, limit, offset)
	}
	serialized := serializers.SerializeSearchResults(searchResults)
	paginateResponse(w, r, serialized, totalCount, limit, offset)
}
