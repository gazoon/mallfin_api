package handlers

import (
	"net/http"

	"mallfin_api/db"
	"mallfin_api/models"
	"mallfin_api/serializers"

	"mallfin_api/logging"

	"github.com/gazoon/binding"
	"github.com/gazoon/httprouter"
)

func mallsBySubwayStation(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	subwayStationID := *formData.SubwayStation
	if !checkSubwayStation(ctx, w, subwayStationID) {
		return
	}
	malls, err := db.GetMallsBySubwayStation(subwayStationID, formData.Sort, formData.Limit, formData.Offset)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return
	}
	totalCount, ok := totalCountFromResults(len(malls), formData.Limit, formData.Offset)
	if !ok {
		logger.Info("Getting count of malls by station from db")
		totalCount, err = db.MallsBySubwayStationCount(subwayStationID)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(ctx, w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func mallsByQuery(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	name := *formData.Query
	var malls []*models.Mall
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		malls, err = db.GetMallsByName(name, userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByNameCount(name, userCity)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		malls, err = db.GetMallsByNameWithoutCity(name, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByNameWithoutCityCount(name)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(ctx, w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func mallsByShop(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	shopID := *formData.Shop
	if !checkShop(ctx, w, shopID) {
		return
	}
	var malls []*models.Mall
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		malls, err = db.GetMallsByShop(shopID, userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByShopCount(shopID, userCity)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		malls, err = db.GetMallsByShopWithoutCity(shopID, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsByShopWithoutCityCount(shopID)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(ctx, w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func allMalls(w http.ResponseWriter, r *http.Request, formData *mallsListForm) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	var malls []*models.Mall
	var totalCount int
	if formData.City != nil {
		userCity := *formData.City
		var err error
		malls, err = db.GetMalls(userCity, formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsCount(userCity)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		var err error
		malls, err = db.GetMallsWithoutCity(formData.Sort, formData.Limit, formData.Offset)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(malls), formData.Limit, formData.Offset)
		if !ok {
			totalCount, err = db.MallsWithoutCityCount()
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeMalls(malls)
	paginateResponse(ctx, w, r, serialized, totalCount, formData.Limit, formData.Offset)
}

func MallsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	formData := &mallsListForm{}
	errs := binding.Form(r, formData)
	if errs != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}

	if !checkCity(ctx, w, formData.City) {
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
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	mallID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	mall, err := db.GetMallDetails(mallID)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return
	}
	if mall == nil {
		notFoundResponse(ctx, w, MALL_NOT_FOUND)
		return
	}
	serialized := serializers.SerializeMall(mall)
	response(ctx, w, serialized)
}

func ShopsInMalls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	formData := shopsInMallsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	mallIDs := formData.Malls
	shopIDs := formData.Shops
	mallsShops, err := db.GetShopsInMalls(mallIDs, shopIDs)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return
	}
	serialized := serializers.SerializeShopsInMalls(mallsShops)
	response(ctx, w, serialized)
}

func CurrentMall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	formData := CoordinatesForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	userLocation := &models.Location{
		Lat: formData.LocationLat,
		Lon: formData.LocationLon,
	}
	mall, err := db.GetMallByLocation(userLocation)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return
	}
	if mall == nil {
		notFoundResponse(ctx, w, MALL_NOT_FOUND)
		return
	}
	serialized := serializers.SerializeMall(mall)
	response(ctx, w, serialized)
}
