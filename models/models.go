package models

import (
	"github.com/pkg/errors"
	"strings"
)

type WeekTime struct {
	Time string
	Day  int
}

type WorkPeriod struct {
	Open  WeekTime
	Close WeekTime
}

type Location struct {
	Lat float64
	Lon float64
}

type Logo struct {
	Small string
	Large string
}

type SubwayStation struct {
	ID   int
	Name string
}

type Mall struct {
	ID         int
	Name       string
	Phone      string
	Logo       Logo
	Location   Location
	ShopsCount int
	Address    string
	//Details
	Site         string
	DayAndNight  bool
	Subway       *SubwayStation
	WorkingHours []*WorkPeriod
}

type Shop struct {
	ID         int
	Name       string
	Logo       Logo
	Score      int
	MallsCount int
	//Details
	Phone       string
	Site        string
	NearestMall *Mall
}

type Category struct {
	ID         int
	Name       string
	Logo       Logo
	ShopsCount int
}

type City struct {
	ID   int
	Name string
}

type MallMatchedShops struct {
	MallID  int
	ShopIDs []int
}

type SearchResult struct {
	Mall     *Mall
	ShopIDs  []int
	Distance *float64
}

type Sorting interface {
	Key() string
	Reversed() bool
}

type sorting struct {
	key      string
	reversed bool
}

func (s *sorting) Key() string {
	return s.key
}

func (s *sorting) Reversed() bool {
	return s.reversed
}

const (
	ID_SORT_KEY          = "id"
	NAME_SORT_KEY        = "name"
	SHOPS_COUNT_SORT_KEY = "shops_count"
	MALLS_COUNT_SORT_KEY = "malls_count"
	SCORE_SORT_KEY       = "score"
	DISTANCE_SORT_KEY    = "distance"
	MALL_NAME_SORT_KEY   = "mall_name"
	MALL_ID_SORT_KEY     = "mall_id"
)

const REVERSE_SIGN = "-"

var (
	DefaultMallSorting     = &sorting{key: ID_SORT_KEY, reversed: false}
	DefaultShopSorting     = DefaultMallSorting
	DefaultCategorySorting = DefaultMallSorting
	DefaultCitySorting     = DefaultMallSorting
	DefaultSearchSorting   = &sorting{key: MALL_ID_SORT_KEY, reversed: false}
)

func MallSorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, ID_SORT_KEY, NAME_SORT_KEY, SHOPS_COUNT_SORT_KEY)
}

func ShopSorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, ID_SORT_KEY, NAME_SORT_KEY, SCORE_SORT_KEY, MALLS_COUNT_SORT_KEY)
}

func CategorySorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, ID_SORT_KEY, NAME_SORT_KEY, SHOPS_COUNT_SORT_KEY)
}

func CitySorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, ID_SORT_KEY, NAME_SORT_KEY)
}

func SearchSorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, MALL_ID_SORT_KEY, MALL_NAME_SORT_KEY, SHOPS_COUNT_SORT_KEY, DISTANCE_SORT_KEY)
}

func modelSorting(rawSorting string, validSortKeys ...string) (Sorting, error) {
	isReversed := strings.HasPrefix(rawSorting, REVERSE_SIGN)
	sortKey := strings.TrimPrefix(rawSorting, REVERSE_SIGN)
	for _, validSortKey := range validSortKeys {
		if sortKey == validSortKey {
			return &sorting{key: sortKey, reversed: isReversed}, nil
		}
	}
	return nil, errors.Errorf("Unsupported sort key: %s, valid values: %v", sortKey, validSortKeys)
}
