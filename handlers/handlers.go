package handlers

import (
	"net/http"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/binding"
	"github.com/gazoon/httprouter"
	"mallfin_api/db/models"
	"mallfin_api/serializers"
)

type mallsListForm struct {
	City          *int
	Shop          *int
	Query         *string
	SubwayStation *int
	Ids           []int
	Sort          *string
	Limit         *uint
	Offset        *uint
}

func (mf *mallsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&mf.City:          "city",
		&mf.Shop:          "shop",
		&mf.SubwayStation: "subway_station",
		&mf.Query:         "query",
		&mf.Ids:           "ids",
		&mf.Sort:          "sort",
		&mf.Limit:         "limit",
		&mf.Offset:        "offset",
	}
}

func (mf *mallsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if mf.Sort != nil {
		sortKey := *mf.Sort
		if sortKey != models.NAME_MALL_SORT_KEY && sortKey != models.SHOPS_COUNT_MALL_SORT_KEY {
			errs = append(errs, binding.Error{
				FieldNames: []string{"sort"},
				Message: fmt.Sprintf("Invalid sort key for list of malls, valid values: %s or %s.",
					models.NAME_MALL_SORT_KEY, models.SHOPS_COUNT_MALL_SORT_KEY),
			})
		}
	}
	return errs
}

func MallsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := mallsListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	//var malls []*models.Mall
	//if formData.Ids != nil {
	//	mallIDs := formData.Ids
	//	malls = models.GetMallsByIds(mallIDs)
	//} else if formData.SubwayStation != nil {
	//	subwayStationID := *formData.SubwayStation
	//	if !models.IsSubwayStationExists(subwayStationID) {
	//		errorResponse(w, SUBWAY_STATION_NOT_FOUND, "Subway station with such id does not exists.", http.StatusNotFound)
	//		return
	//	}
	//	malls = models.GetMallsBySubwayStation(subwayStationID)
	//} else if formData.Query != nil {
	//	name := *formData.Query
	//	if formData.City != nil {
	//		cityID := *formData.City
	//		malls = models.GetMallsByNameAndCity(name, cityID)
	//	} else {
	//		malls = models.GetMallsByName(name)
	//	}
	//} else if formData.Shop != nil {
	//	shopID := *formData.Shop
	//	if formData.City != nil {
	//		cityID := *formData.City
	//		malls = models.GetMallsByShopAndCity(shopID, cityID)
	//	} else {
	//		malls = models.GetMallsByShop(shopID)
	//	}
	//}
	log.Infof("%+v", formData)
	log.Info(formData.Ids == nil)
	log.Info("success")
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
	mallSerialized := serializers.SerializeMallDetails(mall)
	response(w, mallSerialized)
}
