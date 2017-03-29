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
	IDSortKey         = "id"
	NameSortKey       = "name"
	ShopsCountSortKey = "shops_count"
	MallsCountSortKey = "malls_count"
	ScoreSortKey      = "score"
	DistanceSortKey   = "distance"
	MallNameSortKey   = "mall_name"
	MallIDSortKey     = "mall_id"
)

const REVERSE_SIGN = "-"

var (
	DefaultMallSorting     = &sorting{key: IDSortKey, reversed: false}
	DefaultShopSorting     = DefaultMallSorting
	DefaultCategorySorting = DefaultMallSorting
	DefaultCitySorting     = DefaultMallSorting
	DefaultSearchSorting   = &sorting{key: MallIDSortKey, reversed: false}
)

func MallSorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, IDSortKey, NameSortKey, ShopsCountSortKey)
}

func ShopSorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, IDSortKey, NameSortKey, ScoreSortKey, MallsCountSortKey)
}

func CategorySorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, IDSortKey, NameSortKey, ShopsCountSortKey)
}

func CitySorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, IDSortKey, NameSortKey)
}

func SearchSorting(rawSorting string) (Sorting, error) {
	return modelSorting(rawSorting, MallIDSortKey, MallNameSortKey, ShopsCountSortKey, DistanceSortKey)
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
