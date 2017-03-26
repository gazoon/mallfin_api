package models

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"mallfin_api/db"
	"strings"
)

var moduleLog = log.WithField("location", "models")

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
	client := db.GetClient()
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

type SortKeyToOrderBy struct {
	dict           map[string]*OrderBy
	defaultOrderBy *OrderBy
}

var (
	MALL_DEFAULT_ORDER_BY     = &OrderBy{Column: "m.mall_id", Reverse: false}
	SHOP_DEFAULT_ORDER_BY     = &OrderBy{Column: "s.shop_id", Reverse: false}
	CATEGORY_DEFAULT_ORDER_BY = &OrderBy{Column: "c.category_id", Reverse: false}
	CITY_DEFAULT_ORDER_BY     = &OrderBy{Column: "c.city_id", Reverse: false}
	SEARCH_DEFAULT_ORDER_BY   = &OrderBy{Column: "m.mall_id", Reverse: false}

	MALLS_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           MALL_DEFAULT_ORDER_BY,
			"-id":          {Column: "m.mall_id", Reverse: true},
			"name":         {Column: "m.mall_name", Reverse: false},
			"-name":        {Column: "m.mall_name", Reverse: true},
			"shops_count":  {Column: "m.shops_count", Reverse: false},
			"-shops_count": {Column: "m.shops_count", Reverse: true},
		},
		defaultOrderBy: MALL_DEFAULT_ORDER_BY,
	}
	SHOPS_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           SHOP_DEFAULT_ORDER_BY,
			"-id":          {Column: "s.shop_id", Reverse: true},
			"name":         {Column: "s.shop_name", Reverse: false},
			"-name":        {Column: "s.shop_name", Reverse: true},
			"score":        {Column: "s.score", Reverse: false},
			"-score":       {Column: "s.score", Reverse: true},
			"malls_count":  {Column: "s.malls_count", Reverse: false},
			"-malls_count": {Column: "s.malls_count", Reverse: true},
		},
		defaultOrderBy: SHOP_DEFAULT_ORDER_BY,
	}
	CATEGORIES_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":           CATEGORY_DEFAULT_ORDER_BY,
			"-id":          {Column: "c.cateogry_id", Reverse: true},
			"name":         {Column: "c.category_name", Reverse: false},
			"-name":        {Column: "c.category_name", Reverse: true},
			"shops_count":  {Column: "c.shops_count", Reverse: false},
			"-shops_count": {Column: "c.shops_count", Reverse: true},
		},
		defaultOrderBy: CATEGORY_DEFAULT_ORDER_BY,
	}
	CITIES_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"id":    CITY_DEFAULT_ORDER_BY,
			"-id":   {Column: "c.city_id", Reverse: true},
			"name":  {Column: "c.city_name", Reverse: false},
			"-name": {Column: "c.city_name", Reverse: true},
		},
		defaultOrderBy: CITY_DEFAULT_ORDER_BY,
	}
	SEARCH_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"mall_id":      SEARCH_DEFAULT_ORDER_BY,
			"-mall_id":     {Column: "m.mall_id", Reverse: true},
			"mall_name":    {Column: "m.mall_name", Reverse: false},
			"-mall_name":   {Column: "m.mall_name", Reverse: true},
			"shops_count":  {Column: "m.shops_count", Reverse: false},
			"-shops_count": {Column: "m.shops_count", Reverse: true},
		},
		defaultOrderBy: SEARCH_DEFAULT_ORDER_BY,
	}
	SEARCH_WITH_DISTANCE_SORT_KEYS = &SortKeyToOrderBy{
		dict: map[string]*OrderBy{
			"mall_id":      SEARCH_DEFAULT_ORDER_BY,
			"-mall_id":     {Column: "m.mall_id", Reverse: true},
			"mall_name":    {Column: "m.mall_name", Reverse: false},
			"-mall_name":   {Column: "m.mall_name", Reverse: true},
			"shops_count":  {Column: "m.shops_count", Reverse: false},
			"-shops_count": {Column: "m.shops_count", Reverse: true},
			"distance":     {Column: "distance", Reverse: false},
			"-distance":    {Column: "distance", Reverse: true},
		},
		defaultOrderBy: SEARCH_DEFAULT_ORDER_BY,
	}
)

func (sk *SortKeyToOrderBy) FmtKeys() string {
	keys := make([]string, 0, len(sk.dict))
	for key := range sk.dict {
		keys = append(keys, key)
	}
	return strings.Join(keys, ", ")
}

func (sk *SortKeyToOrderBy) IsValid(sortKey *string) bool {
	if sortKey != nil {
		if _, ok := sk.dict[*sortKey]; !ok {
			return false
		}
	}
	return true
}

func (sk *SortKeyToOrderBy) CorrespondingOrderBy(sortKey *string) *OrderBy {
	orderBy := sk.defaultOrderBy
	if sortKey != nil {
		if correspondOrderBy, ok := sk.dict[*sortKey]; ok {
			orderBy = correspondOrderBy
		}
	}
	return orderBy
}
