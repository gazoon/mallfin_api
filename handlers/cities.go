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

func CitiesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := citiesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	var cities []*models.City
	if formData.Query != nil {
		name := *formData.Query
		var err error
		cities, err = db.GetCitiesByName(name, sorting)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
	} else {
		var err error
		cities, err = db.GetCities(sorting)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
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
