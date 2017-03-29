package handlers

import (
	"mallfin_api/models"
	"net/http"

	"github.com/gazoon/binding"
)

type checkSortKeyFn func(string) (models.Sorting, error)

func bindSortKey(toSorting checkSortKeyFn, fieldName string, formVals []string, errs *binding.Errors) models.Sorting {
	if len(formVals) == 0 {
		return nil
	}
	sortKey := formVals[0]
	sorting, err := toSorting(sortKey)
	if err != nil {
		errs.Add([]string{fieldName}, "", err.Error())
	}
	return sorting
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

type mallsListForm struct {
	City          *int
	Shop          *int
	Query         *string
	SubwayStation *int
	Sort          models.Sorting
	Limit         *int
	Offset        *int
}

func (mlf *mallsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&mlf.City:          "city",
		&mlf.Shop:          "shop",
		&mlf.SubwayStation: "subway_station",
		&mlf.Query:         "query",
		&mlf.Limit:         "limit",
		&mlf.Offset:        "offset",
		&mlf.Sort: binding.Field{
			Form: "sort",
			Binder: func(fieldName string, formVals []string, errs binding.Errors) binding.Errors {
				sorting := bindSortKey(models.MallSorting, fieldName, formVals, &errs)
				mlf.Sort = sorting
				return errs
			},
		},
	}
}

func (mlf *mallsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	errs = checkLimitOffset(mlf.Limit, mlf.Offset, errs)
	return errs
}

type shopsListForm struct {
	City     *int
	Mall     *int
	Query    *string
	Category *int
	Sort     models.Sorting
	Limit    *int
	Offset   *int
}

func (slf *shopsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&slf.City:     "city",
		&slf.Mall:     "mall",
		&slf.Category: "category",
		&slf.Query:    "query",
		&slf.Limit:    "limit",
		&slf.Offset:   "offset",
		&slf.Sort: binding.Field{
			Form: "sort",
			Binder: func(fieldName string, formVals []string, errs binding.Errors) binding.Errors {
				sorting := bindSortKey(models.ShopSorting, fieldName, formVals, &errs)
				slf.Sort = sorting
				return errs
			},
		},
	}
}

func (slf *shopsListForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
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
	Sort models.Sorting
}

func (clf *categoriesListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&clf.City: "city",
		&clf.Shop: "shop",
		&clf.Sort: binding.Field{
			Form: "sort",
			Binder: func(fieldName string, formVals []string, errs binding.Errors) binding.Errors {
				sorting := bindSortKey(models.CategorySorting, fieldName, formVals, &errs)
				clf.Sort = sorting
				return errs
			},
		},
	}
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
	Sort  models.Sorting
}

func (clf *citiesListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&clf.Query: "query",
		&clf.Sort: binding.Field{
			Form: "sort",
			Binder: func(fieldName string, formVals []string, errs binding.Errors) binding.Errors {
				sorting := bindSortKey(models.CitySorting, fieldName, formVals, &errs)
				clf.Sort = sorting
				return errs
			},
		},
	}
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
	Sort        models.Sorting
	Limit       *int
	Offset      *int
}

func (sf *searchForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&sf.Shops:       "shops",
		&sf.City:        "city",
		&sf.LocationLat: "location_lat",
		&sf.LocationLon: "location_lon",
		&sf.Limit:       "limit",
		&sf.Offset:      "offset",
		&sf.Sort: binding.Field{
			Form: "sort",
			Binder: func(fieldName string, formVals []string, errs binding.Errors) binding.Errors {
				sorting := bindSortKey(models.SearchSorting, fieldName, formVals, &errs)
				sf.Sort = sorting
				return errs
			},
		},
	}
}

func (sf *searchForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	sorting := sf.Sort
	if sorting != nil && sorting.Key() == models.DistanceSortKey && (sf.LocationLon == nil || sf.LocationLat == nil) {
		errs = append(errs, binding.Error{
			FieldNames: []string{"sort"},
			Message:    "cannot sort by distance without location",
		})
		return errs
	}
	errs = checkLimitOffset(sf.Limit, sf.Offset, errs)
	return errs
}
