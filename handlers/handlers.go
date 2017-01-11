package handlers

import (
	"net/http"

	"fmt"
	"mallfin_api/db/models"
	"mallfin_api/serializers"

	"github.com/gazoon/binding"
	"github.com/gazoon/httprouter"
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
	if !models.MALL_SORT_KEYS.IsValid(mf.Sort) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of malls, valid values: %s.", models.MALL_SORT_KEYS.FmtKeys()),
		})
	}
	return errs
}

type shopsListForm struct {
	City     *int
	Mall     *int
	Query    *string
	Category *int
	Ids      []int
	Sort     *string
	Limit    *uint
	Offset   *uint
}

func (sf *shopsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&sf.City:     "city",
		&sf.Mall:     "mall",
		&sf.Category: "category",
		&sf.Query:    "query",
		&sf.Ids:      "ids",
		&sf.Sort:     "sort",
		&sf.Limit:    "limit",
		&sf.Offset:   "offset",
	}
}

func (sf *shopsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if !models.SHOP_SORT_KEYS.IsValid(sf.Sort) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of shops, valid values: %s.", models.SHOP_SORT_KEYS.FmtKeys()),
		})
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
	mallsSerialized := serializers.SerializeMalls(malls)
	listResponse(w, mallsSerialized, totalCount)
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
	mallSerialized := serializers.SerializeMall(mall)
	objectResponse(w, mallSerialized)
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
			errorResponse(w, MALL_NOT_FOUND, "mall with such id does not exists.", http.StatusNotFound)
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
	shopsSerialized := serializers.SerializeShops(shops)
	listResponse(w, shopsSerialized, totalCount)

}
func ShopDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shopID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INVALID_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	shop := models.GetShopDetails(shopID)
	if shop == nil {
		errorResponse(w, SHOP_NOT_FOUND, "Shop with such id does not exists", http.StatusNotFound)
		return
	}
	shopSerialized := serializers.SerializeShop(shop)
	objectResponse(w, shopSerialized)
}
