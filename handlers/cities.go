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

func CitiesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	formData := citiesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	var cities []*models.City
	if formData.Query != nil {
		name := *formData.Query
		var err error
		cities, err = db.GetCitiesByName(name, sorting)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
	} else {
		var err error
		cities, err = db.GetCities(sorting)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
	}
	serialized := serializers.SerializeCities(cities)
	response(ctx, w, serialized)
}

func CurrentCity(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	city, err := db.GetCityByLocation(userLocation)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return
	}
	if city == nil {
		errorResponse(ctx, w, CITY_NOT_FOUND, "In this place there is no city.", http.StatusNotFound)
		return
	}
	serialized := serializers.SerializeCity(city)
	response(ctx, w, serialized)
}
