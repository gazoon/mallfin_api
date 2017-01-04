package main

import (
	//json "github.com/pquerna/ffjson/ffjson"
	"encoding/json"
	"testing"
)

type TestModel struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Address       string         `json:"address"`
	Site          string         `json:"siteee"`
	Phone         string         `json:"phone"`
	DayAndNight   bool           `json:"day_and_night"`
	CityID        int            `json:"-"`
	Location      *Location      `json:"location"`
	SubwayStation *SubwayStation `json:"subway_station"`
	WorkingHours  []*WorkPeriod  `json:"working_hours"`
}
type SubwayStation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type WorkPeriod struct {
	Open  int `json:"opening"`
	Close int `json:"closing"`
}

var marshaledData = struct {
	ID            int
	Name          string
	Address       string
	Site          string
	Phone         string
	DayAndNight   bool
	CityID        int
	Location      *Location
	SubwayStation *SubwayStation
	WorkingHours  []*WorkPeriod
}{
	ID:          1,
	Name:        "aaaaa",
	Address:     "dsfasfdsfsadfasdfasdf",
	Site:        "asdfjas;djfsldjflsdjfsdf",
	Phone:       "adsfsadf",
	DayAndNight: false,
	CityID:      22,
	Location: &Location{
		Lat: 222.3333,
		Lon: 222.3333,
	},
	SubwayStation: &SubwayStation{
		ID:   228,
		Name: "asdfasdf",
	},
	WorkingHours: []*WorkPeriod{
		{Open: 1, Close: 1},
		{Open: 2, Close: 2},
		{Open: 3, Close: 3},
		{Open: 4, Close: 4},
		{Open: 5, Close: 5},
		{Open: 6, Close: 6},
		{Open: 7, Close: 7},
	},
}

func marshalStruct() {
	t := TestModel{
		ID:            marshaledData.ID,
		Name:          marshaledData.Name,
		Address:       marshaledData.Address,
		Site:          marshaledData.Site,
		Phone:         marshaledData.Phone,
		DayAndNight:   marshaledData.DayAndNight,
		CityID:        marshaledData.CityID,
		Location:      marshaledData.Location,
		SubwayStation: marshaledData.SubwayStation,
		WorkingHours:  marshaledData.WorkingHours,
	}
	_, err := json.Marshal(&t)
	if err != nil {
		panic(err)
	}
}

type JSON map[string]interface{}

func marshalMap() {
	t := JSON{
		"id":             marshaledData.ID,
		"name":           marshaledData.Name,
		"address":        marshaledData.Address,
		"site":           marshaledData.Site,
		"phone":          marshaledData.Phone,
		"day_and_night":  marshaledData.DayAndNight,
		"city_id":        marshaledData.CityID,
		"location":       marshaledData.Location,
		"subway_station": marshaledData.SubwayStation,
		"working_hours":  marshaledData.WorkingHours,
	}
	_, err := json.Marshal(&t)
	if err != nil {
		panic(err)
	}
}
func BenchmarkMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		marshalMap()
	}
}
func BenchmarkStruct(b *testing.B) {
	for n := 0; n < b.N; n++ {
		marshalStruct()
	}
}
