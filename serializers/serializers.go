package serializers

import (
	"mallfin_api/db/models"
	"time"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type Logo struct {
	Large string `json:"large"`
	Small string `json:"small"`
}
type SubwayStation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type WeekTime struct {
	Time time.Time `json:"time"`
	Day  int       `json:"day"`
}
type WorkPeriod struct {
	Opening *WeekTime `json:"opening"`
	Closing *WeekTime `json:"closing"`
}
type Mall struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Phone         string         `json:"phone"`
	Address       string         `json:"address"`
	Logo          *Logo          `json:"logo"`
	Location      *Location      `json:"location"`
	SubwayStation *SubwayStation `json:"subway_staion"`
}
type MallDetails struct {
	*Mall
	Site         string        `json:"site"`
	DayAndNight  bool          `json:"day_and_night"`
	WorkingHours []*WorkPeriod `json:"working_hours"`
}

func SerializeMall(mall *models.Mall) *Mall {
	var subwayStation *SubwayStation
	if mall.SubwayID != nil && mall.SubwayName != nil {
		subwayStation = &SubwayStation{ID: *mall.SubwayID, Name: *mall.SubwayName}
	}
	ms := &Mall{
		ID:      mall.ID,
		Name:    mall.Name,
		Phone:   mall.Phone,
		Address: mall.Address,
		Logo: &Logo{
			Large: mall.LogoLarge,
			Small: mall.LogoSmall,
		},
		Location: &Location{
			Lat: mall.LocationLat,
			Lon: mall.LocationLon,
		},
		SubwayStation: subwayStation,
	}
	return ms
}
func SerializeMallDetails(mallDetails *models.MallDetails) *MallDetails {
	workingHours := []*WorkPeriod{}
	for _, period := range mallDetails.WorkingHours {
		workingHours = append(workingHours, &WorkPeriod{
			Opening: &WeekTime{
				Time: period.OpenTime,
				Day:  period.OpenDay,
			},
			Closing: &WeekTime{
				Time: period.CloseTime,
				Day:  period.CloseDay,
			},
		})
	}
	mds := &MallDetails{
		Mall:         SerializeMall(mallDetails.Mall),
		Site:         mallDetails.Site,
		DayAndNight:  mallDetails.DayAndNight,
		WorkingHours: workingHours,
	}
	return mds
}
