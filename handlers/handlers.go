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
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sortKey := formData.Sort
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	if cityID != nil {
		if !models.IsCityExists(*cityID) {
			errorResponse(w, CITY_NOT_FOUND, "City with such id does not exists.", http.StatusNotFound)
			return
		}
	}
	var malls []*models.Mall
	var totalCount int
	if formData.Ids != nil {
		mallIDs := formData.Ids
		malls, totalCount = models.GetMallsByIds(mallIDs)
	} else if formData.SubwayStation != nil {
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
	listResponse(w, serialized, totalCount)
}
func MallDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mallID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INVALID_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	mall := models.GetMallDetails(mallID)
	if mall == nil {
		errorResponse(w, MALL_NOT_FOUND, "Mall with such id does not exists", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeMall(mall)
	objectResponse(w, serialized)
}
func ShopsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := shopsListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sortKey := formData.Sort
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	if cityID != nil {
		if !models.IsCityExists(*cityID) {
			errorResponse(w, CITY_NOT_FOUND, "City with such id does not exists.", http.StatusNotFound)
			return
		}
	}
	var shops []*models.Shop
	var totalCount int
	if formData.Ids != nil {
		shopIDs := formData.Ids
		shops, totalCount = models.GetShopsByIds(shopIDs, cityID)
	} else if formData.Mall != nil {
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
	listResponse(w, serialized, totalCount)
}
func ShopDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	formData := shopDetailsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	shopID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INVALID_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	cityID := formData.City
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
	objectResponse(w, serialized)
}
func CategoriesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := categoriesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sortKey := formData.Sort
	cityID := formData.City
	if cityID != nil {
		if !models.IsCityExists(*cityID) {
			errorResponse(w, CITY_NOT_FOUND, "City with such id does not exists.", http.StatusNotFound)
			return
		}
	}
	var categories []*models.Category
	var totalCount int
	if formData.Ids != nil {
		categoryIDs := formData.Ids
		categories, totalCount = models.GetCategoriesByIds(categoryIDs, cityID)
	} else if formData.Shop != nil {
		shopID := *formData.Shop
		if !models.IsShopExists(shopID) {
			errorResponse(w, SHOP_NOT_FOUND, "Shop with such id does not exists.", http.StatusNotFound)
			return
		}
		categories, totalCount = models.GetCategoriesByShop(shopID, cityID, sortKey)
	} else {
		categories, totalCount = models.GetCategories(cityID, sortKey)
	}
	serialized := serializers.SerializeCategories(categories)
	listResponse(w, serialized, totalCount)
}
func CategoryDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	formData := categoryDetailsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	categoryID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INVALID_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	cityID := formData.City
	category := models.GetCategoryDetails(categoryID, cityID)
	if category == nil {
		errorResponse(w, CATEGORY_NOT_FOUND, "Category with such id does not exists", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeCategory(category)
	objectResponse(w, serialized)
}
func CitiesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := citiesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sortKey := formData.Sort
	var cities []*models.City
	var totalCount int
	if formData.Query != nil {
		name := *formData.Query
		cities, totalCount = models.GetCitiesByName(name, sortKey)
	} else {
		cities, totalCount = models.GetCities(sortKey)
	}
	serialized := serializers.SerializeCities(cities)
	listResponse(w, serialized, totalCount)
}
func CurrentMall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := currentMallForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
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
	objectResponse(w, serialized)
}
func ShopsInMalls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := shopsInMallsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	mallIDs := formData.Malls
	shopIDs := formData.Shops
	mallsShops := models.GetShopsInMalls(mallIDs, shopIDs)
	serialized := serializers.SerializeShopsInMalls(mallsShops)
	objectResponse(w, serialized)
}
