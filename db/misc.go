package db

import (
	"fmt"
	"strings"

	"mallfin_api/models"

	"github.com/pkg/errors"
)

type baseQuery string

func (bq baseQuery) withColumns(columns string) string {
	return strings.Replace(string(bq), "{columns}", columns, 1)
}

func totalCountFromResults(resultsLen int, limit, offset *int) (int, bool) {
	if (limit == nil || *limit == 0) && (offset == nil || *offset == 0 || resultsLen != 0) {
		totalCount := resultsLen
		if offset != nil {
			totalCount += *offset
		}
		return totalCount, true
	}
	return 0, false
}

func countQuery(queryName, query string, args ...interface{}) (int, error) {
	var row struct {
		Count int
	}
	client := GetClient()
	_, err := client.QueryOne(&row, query, args...)
	if err != nil {
		return 0, errors.WithMessage(err, queryName)
	}
	return row.Count, nil
}

type OrderBy struct {
	Column  string
	Reverse bool
}

func (o *OrderBy) String() string {
	return o.ToSql()
}

func (o *OrderBy) ToSql() string {
	var s string
	if o.Reverse {
		s = fmt.Sprintf("%s DESC", o.Column)
	} else {
		s = fmt.Sprintf("%s ASC", o.Column)
	}
	return s
}

func (o *OrderBy) CompileQuery(query string) string {
	return strings.Replace(query, "{order}", o.ToSql(), 1)
}

func (o *OrderBy) CompileBaseQuery(query string) baseQuery {
	return baseQuery(o.CompileQuery(query))
}

func mallOrderBy(sorting models.Sorting) *OrderBy {
	if sorting == nil {
		sorting = models.DefaultMallSorting
	}
	var column string
	switch sorting.Key() {
	case models.IDSortKey:
		column = "m.mall_id"
	case models.NameSortKey:
		column = "m.mall_name"
	case models.ShopsCountSortKey:
		column = "m.shops_count"
	default:
		panic(errors.Errorf("Unexpected sorting key %s for mall order by", sorting.Key()))
	}
	return &OrderBy{Column: column, Reverse: sorting.Reversed()}
}

func shopOrderBy(sorting models.Sorting) *OrderBy {
	if sorting == nil {
		sorting = models.DefaultShopSorting
	}
	var column string
	switch sorting.Key() {
	case models.IDSortKey:
		column = "s.shop_id"
	case models.NameSortKey:
		column = "s.shop_name"
	case models.MallsCountSortKey:
		column = "s.malls_count"
	case models.ScoreSortKey:
		column = "s.score"
	default:
		panic(errors.Errorf("Unexpected sorting key %s for shop order by", sorting.Key()))
	}
	return &OrderBy{Column: column, Reverse: sorting.Reversed()}
}

func categoryOrderBy(sorting models.Sorting) *OrderBy {
	if sorting == nil {
		sorting = models.DefaultCategorySorting
	}
	var column string
	switch sorting.Key() {
	case models.IDSortKey:
		column = "c.category_id"
	case models.NameSortKey:
		column = "c.category_name"
	case models.ShopsCountSortKey:
		column = "c.shops_count"
	default:
		panic(errors.Errorf("Unexpected sorting key %s for category order by", sorting.Key()))
	}
	return &OrderBy{Column: column, Reverse: sorting.Reversed()}
}

func cityOrderBy(sorting models.Sorting) *OrderBy {
	if sorting == nil {
		sorting = models.DefaultCitySorting
	}
	var column string
	switch sorting.Key() {
	case models.IDSortKey:
		column = "c.city_id"
	case models.NameSortKey:
		column = "c.city_name"
	default:
		panic(errors.Errorf("Unexpected sorting key %s for city order by", sorting.Key()))
	}
	return &OrderBy{Column: column, Reverse: sorting.Reversed()}
}

func searchOrderBy(sorting models.Sorting) *OrderBy {
	if sorting == nil {
		sorting = models.DefaultSearchSorting
	}
	var column string
	switch sorting.Key() {
	case models.MallIDSortKey:
		column = "m.mall_id"
	case models.MallNameSortKey:
		column = "m.mall_name"
	case models.ShopsCountSortKey:
		column = "m.shops_count"
	case models.DistanceSortKey:
		column = "distance"
	default:
		panic(errors.Errorf("Unexpected sorting key %s for search order by", sorting.Key()))
	}
	return &OrderBy{Column: column, Reverse: sorting.Reversed()}
}
