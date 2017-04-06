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

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	formData := searchForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	limit := formData.Limit
	offset := formData.Offset
	cityID := formData.City
	shopIDs := formData.Shops
	sorting := formData.Sort
	if !checkCity(ctx, w, cityID, "log prefix") {
		return
	}
	var searchResults []*models.SearchResult
	var totalCount int
	var userLocation *models.Location
	if formData.LocationLat != nil && formData.LocationLon != nil {
		userLocation = &models.Location{
			Lat: *formData.LocationLat,
			Lon: *formData.LocationLon,
		}
	}
	if cityID != nil {
		userCity := *cityID
		if userLocation != nil {
			var err error
			searchResults, err = db.GetSearchResultsWithDistance(shopIDs, userLocation, userCity, sorting, limit, offset)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		} else {
			var err error
			searchResults, err = db.GetSearchResults(shopIDs, userCity, sorting, limit, offset)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(searchResults), limit, offset)
		if !ok {
			var err error
			totalCount, err = db.SearchResultsCount(shopIDs, userCity)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	} else {
		if userLocation != nil {
			var err error
			searchResults, err = db.GetSearchResultsWithDistanceWithoutCity(shopIDs, userLocation, sorting, limit, offset)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		} else {
			var err error
			searchResults, err = db.GetSearchResultsWithoutCity(shopIDs, sorting, limit, offset)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
		var ok bool
		totalCount, ok = totalCountFromResults(len(searchResults), limit, offset)
		if !ok {
			var err error
			totalCount, err = db.SearchResultsWithoutCityCount(shopIDs)
			if err != nil {
				logger.Error(err)
				internalErrorResponse(w)
				return
			}
		}
	}
	serialized := serializers.SerializeSearchResults(searchResults)
	paginateResponse(ctx, w, r, serialized, totalCount, limit, offset)
}
