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
	sortKeys := models.MALL_SORT_KEYS
	if !sortKeys.IsValid(mf.Sort) {
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
	sortKeys := models.SHOP_SORT_KEYS
	if !sortKeys.IsValid(sf.Sort) {
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

func (cf *shopDetailsForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.City:        "city",
		&cf.LocationLat: "location_lat",
		&cf.LocationLon: "location_lon",
	}
}

type categoriesListForm struct {
	City *int
	Shop *int
	Ids  []int
	Sort *string
}

func (cf *categoriesListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.City: "city",
		&cf.Shop: "shop",
		&cf.Ids:  "ids",
		&cf.Sort: "sort",
	}
}
func (cf *categoriesListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.CATEGORY_SORT_KEYS
	if !sortKeys.IsValid(cf.Sort) {
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

func (cf *categoryDetailsForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.City: "city",
	}
}

type citiesListForm struct {
	Query *string
	Sort  *string
}

func (cf *citiesListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Query: "query",
		&cf.Sort:  "sort",
	}
}
func (cf *citiesListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sortKeys := models.CITY_SORT_KEYS
	if !sortKeys.IsValid(cf.Sort) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    fmt.Sprintf("Invalid sort key for list of cities, valid values: %s.", sortKeys.FmtKeys()),
		})
	}
	return errs
}
