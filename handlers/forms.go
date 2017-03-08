package handlers

import (
	"fmt"
	"github.com/gazoon/binding"
	"mallfin_api/db/models"
	"net/http"
)

type mallsListForm struct {
	City          *int
	Shop          *int
	Query         *string
	SubwayStation *int
	Sort          *string
	Limit         *int
	Offset        *int
}

func (mlf *mallsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&mlf.City:          "city",
		&mlf.Shop:          "shop",
		&mlf.SubwayStation: "subway_station",
		&mlf.Query:         "query",
		&mlf.Sort:          "sort",
		&mlf.Limit:         "limit",
		&mlf.Offset:        "offset",
	}
}

func (mlf *mallsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.MALLS_SORT_KEYS
	errs = checkSortKey(sortKeys, mlf.Sort, errs)
	errs = checkLimitOffset(mlf.Limit, mlf.Offset, errs)
	return errs
}

type shopsListForm struct {
	City     *int
	Mall     *int
	Query    *string
	Category *int
	Sort     *string
	Limit    *int
	Offset   *int
}

func (slf *shopsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&slf.City:     "city",
		&slf.Mall:     "mall",
		&slf.Category: "category",
		&slf.Query:    "query",
		&slf.Sort:     "sort",
		&slf.Limit:    "limit",
		&slf.Offset:   "offset",
	}
}

func (slf *shopsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.SHOPS_SORT_KEYS
	errs = checkSortKey(sortKeys, slf.Sort, errs)
	errs = checkLimitOffset(slf.Limit, slf.Offset, errs)
	return errs
}

type shopDetailsForm struct {
	City        *int
	LocationLat *float64
	LocationLon *float64
}

func (sdf *shopDetailsForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&sdf.City:        "city",
		&sdf.LocationLat: "location_lat",
		&sdf.LocationLon: "location_lon",
	}
}

type categoriesListForm struct {
	City *int
	Shop *int
	Sort *string
}

func (clf *categoriesListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&clf.City: "city",
		&clf.Shop: "shop",
		&clf.Sort: "sort",
	}
}
func (clf *categoriesListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.CATEGORIES_SORT_KEYS
	return checkSortKey(sortKeys, clf.Sort, errs)
}

type categoryDetailsForm struct {
	City *int
}

func (cdf *categoryDetailsForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cdf.City: "city",
	}
}

type citiesListForm struct {
	Query *string
	Sort  *string
}

func (clf *citiesListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&clf.Query: "query",
		&clf.Sort:  "sort",
	}
}
func (clf *citiesListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.CITIES_SORT_KEYS
	return checkSortKey(sortKeys, clf.Sort, errs)
}

type CoordinatesForm struct {
	LocationLat float64
	LocationLon float64
}

func (cmf *CoordinatesForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cmf.LocationLat: binding.Field{
			Form:     "location_lat",
			Required: true,
		},
		&cmf.LocationLon: binding.Field{
			Form:     "location_lon",
			Required: true,
		},
	}
}

type shopsInMallsForm struct {
	Shops []int
	Malls []int
}

func (smf *shopsInMallsForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&smf.Shops: "shops",
		&smf.Malls: "malls",
	}
}

type searchForm struct {
	Shops       []int
	City        *int
	LocationLat *float64
	LocationLon *float64
	Sort        *string
	Limit       *int
	Offset      *int
}

func (sf *searchForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&sf.Shops:       "shops",
		&sf.City:        "city",
		&sf.LocationLat: "location_lat",
		&sf.LocationLon: "location_lon",
		&sf.Sort:        "sort",
		&sf.Limit:       "limit",
		&sf.Offset:      "offset",
	}
}

func (sf *searchForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	var sortKeys *models.SortKeyToOrderBy
	if sf.LocationLat != nil && sf.LocationLon != nil {
		sortKeys = models.SEARCH_WITH_DISTANCE_SORT_KEYS
	} else {
		sortKeys = models.SEARCH_SORT_KEYS
	}
	errs = checkSortKey(sortKeys, sf.Sort, errs)
	errs = checkLimitOffset(sf.Limit, sf.Offset, errs)
	return errs
}

func checkSortKey(validSortKeys *models.SortKeyToOrderBy, sortKey *string, errs binding.Errors) binding.Errors {
	if !validSortKeys.IsValid(sortKey) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of cities, valid values: %s.", validSortKeys.FmtKeys()),
		})
	}
	return errs
}
func checkLimitOffset(limit, offset *int, errs binding.Errors) binding.Errors {
	if limit != nil && *limit < 0 {
		errs = append(errs, binding.Error{
			FieldNames: []string{"limit"},
			Message:    "limit must be non-negative int",
		})
	}
	if offset != nil && *offset < 0 {
		errs = append(errs, binding.Error{
			FieldNames: []string{"offset"},
			Message:    "offset must be non-negative int",
		})
	}
	return errs
}
