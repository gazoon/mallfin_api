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

func shopsByMall(w http.ResponseWriter, r *http.Request, formData *shopsListForm) {
	mallID := *formData.Mall
	if !checkMall(w, mallID, "log prefix") {
		return
	}
	shops, err := db.GetShopsByMall(mallID, formData.Sort, formData.Limit, formData.Offset)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	totalCount, ok := totalCountFromResults(len(shops), formData.Limit, formData.Offset)
	if !ok {
		totalCount, err = db.ShopsByMallCount(mallID)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
	}
	serialized := serializers.SerializeShops(shops)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func shopsByQuery(w http.ResponseWriter, r *http.Request, formData *shopsListForm) {
	name := *formData.Query
	var shops []*models.Shop
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		shops, err = db.GetShopsByName(name, userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(shops), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.ShopsByNameCount(name, userCity)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		shops, err = db.GetShopsByNameWithoutCity(name, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(shops), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.ShopsByNameWithoutCityCount(name)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeShops(shops)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func shopsByCategory(w http.ResponseWriter, r *http.Request, formData *shopsListForm) {
	categoryID := *formData.Category
	if !checkCategory(w,categoryID, "log prefix") {
		return
	}
	var shops []*models.Shop
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		shops, err = db.GetShopsByCategory(categoryID, userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(shops), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.ShopsByCategoryCount(categoryID, userCity)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		shops, err = db.GetShopsByCategoryWithoutCity(categoryID, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(shops), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.ShopsByCategoryWithoutCityCount(categoryID)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeShops(shops)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func allShops(w http.ResponseWriter, r *http.Request, formData *shopsListForm) {
	var shops []*models.Shop
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		shops, err = db.GetShops(userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(shops), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.ShopsCount(userCity)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		shops, err = db.GetShopsWithoutCity(formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(shops), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.ShopsWithoutCityCount()
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeShops(shops)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func ShopsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := &shopsListForm{}
	errs := binding.Form(r, formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}

	if !checkCity(w, formData.City, "log prefix") {
		return
	}

	if formData.Mall != nil {
		shopsByMall(w, r, formData)
	} else if formData.Query != nil {
		shopsByQuery(w, r, formData)
	} else if formData.Category != nil {
		shopsByCategory(w, r, formData)
	} else {
		allShops(w, r, formData)
	}
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

	if !checkCity(w, formData.City, "log prefix") {
		return
	}

	var shop *models.Shop
	if formData.LocationLat != nil && formData.LocationLon != nil {
		userLocation := &models.Location{
			Lat: *formData.LocationLat,
			Lon: *formData.LocationLon,
		}
		shop, err = db.GetShopDetailsWithLocation(shopID, userLocation)
	} else {
		shop, err = db.GetShopDetails(shopID)
	}
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	if shop == nil {
		notFoundResponse(w, SHOP_NOT_FOUND)
		return
	}
	serialized := serializers.SerializeShop(shop)
	response(w, serialized)
}
