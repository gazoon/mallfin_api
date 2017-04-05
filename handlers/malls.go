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

func mallsBySubwayStation(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	subwayStationID := *formData.SubwayStation
	if !checkSubwayStation(w, subwayStationID, "log prefix") {
		return
	}
	malls, err := db.GetMallsBySubwayStation(subwayStationID, formData.Sort, formData.Limit, formData.Offset)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	totalCount, ok := totalCountFromResults(len(malls), formData.Limit, formData.Offset)
	if !ok {
		totalCount, err = db.MallsBySubwayStationCount(subwayStationID)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func mallsByQuery(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	name := *formData.Query
	var malls []*models.Mall
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		malls, err = db.GetMallsByName(name, userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByNameCount(name, userCity)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		malls, err = db.GetMallsByNameWithoutCity(name, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByNameWithoutCityCount(name)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func mallsByShop(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	shopID := *formData.Shop
	if !checkShop(w, shopID, "log prefix") {
		return
	}
	var malls []*models.Mall
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		malls, err = db.GetMallsByShop(shopID, userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByShopCount(shopID, userCity)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		malls, err = db.GetMallsByShopWithoutCity(shopID, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByShopWithoutCityCount(shopID)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func allMalls(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	var malls []*models.Mall
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		malls, err = db.GetMalls(userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsCount(userCity)
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		malls, err = db.GetMallsWithoutCity(formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsWithoutCityCount()
			if err != nil {
				log.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func MallsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := &mallsListForm{}
	errs := binding.Form(r, formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}

	if !checkCity(w, formData.City, "log prefix") {
		return
	}

	if formData.SubwayStation != nil {
		mallsBySubwayStation(w, r, formData)
	} else if formData.Query != nil {
		mallsByQuery(w, r, formData)
	} else if formData.Shop != nil {
		mallsByShop(w, r, formData)
	} else {
		allMalls(w, r, formData)
	}
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
		notFoundResponse(w, MALL_NOT_FOUND)
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
		notFoundResponse(w, MALL_NOT_FOUND)
		return
	}
	serialized := serializers.SerializeMall(mall)
	response(w, serialized)
}
