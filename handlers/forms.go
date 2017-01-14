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
	Ids           []int
	Sort          *string
	Limit         *uint
	Offset        *uint
}

func (mlf *mallsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&mlf.City:          "city",
		&mlf.Shop:          "shop",
		&mlf.SubwayStation: "subway_station",
		&mlf.Query:         "query",
		&mlf.Ids:           "ids",
		&mlf.Sort:          "sort",
		&mlf.Limit:         "limit",
		&mlf.Offset:        "offset",
	}
}

func (mlf *mallsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.MALLS_SORT_KEYS
	if !sortKeys.IsValid(mlf.Sort) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of malls, valid values: %s.", sortKeys.FmtKeys()),
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

func (slf *shopsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&slf.City:     "city",
		&slf.Mall:     "mall",
		&slf.Category: "category",
		&slf.Query:    "query",
		&slf.Ids:      "ids",
		&slf.Sort:     "sort",
		&slf.Limit:    "limit",
		&slf.Offset:   "offset",
	}
}

func (slf *shopsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.SHOPS_SORT_KEYS
	if !sortKeys.IsValid(slf.Sort) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of shops, valid values: %s.", sortKeys.FmtKeys()),
		})
	}
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
	Ids  []int
	Sort *string
}

func (clf *categoriesListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&clf.City: "city",
		&clf.Shop: "shop",
		&clf.Ids:  "ids",
		&clf.Sort: "sort",
	}
}
func (clf *categoriesListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.CATEGORIES_SORT_KEYS
	if !sortKeys.IsValid(clf.Sort) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of categories, valid values: %s.", sortKeys.FmtKeys()),
		})
	}
	return errs
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
	if !sortKeys.IsValid(clf.Sort) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of cities, valid values: %s.", sortKeys.FmtKeys()),
		})
	}
	return errs
}

type currentMallForm struct {
	LocationLat float64
	LocationLon float64
}

func (cmf *currentMallForm) FieldMap(req *http.Request) binding.FieldMap {
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
	Limit       *uint
	Offset      *uint
}

func (sf *searchForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&sf.Shops:       "shops",
		&sf.City:        "city",
		&sf.LocationLat: "location_lat",
		&sf.LocationLon: "location_lon",
		&sf.Limit:       "limit",
		&sf.Offset:      "offset",
	}
}
