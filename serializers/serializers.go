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
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Logo     *Logo     `json:"logo"`
	Location *Location `json:"location"`
}
type MallDetails struct {
	*Mall
	Address       string         `json:"address"`
	Site          string         `json:"site"`
	DayAndNight   bool           `json:"day_and_night"`
	WorkingHours  []*WorkPeriod  `json:"working_hours"`
	SubwayStation *SubwayStation `json:"subway_staion"`
}

func SerializeMall(mall *models.Mall) *Mall {
	ms := &Mall{
		ID:    mall.ID,
		Name:  mall.Name,
		Phone: mall.Phone,
		Logo: &Logo{
			Large: mall.LogoLarge,
			Small: mall.LogoSmall,
		},
		Location: &Location{
			Lat: mall.LocationLat,
			Lon: mall.LocationLon,
		},
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
	var subwayStation *SubwayStation
	if mallDetails.SubwayID != nil && mallDetails.SubwayName != nil {
		subwayStation = &SubwayStation{ID: *mallDetails.SubwayID, Name: *mallDetails.SubwayName}
	}
	mds := &MallDetails{
		Mall:          SerializeMall(mallDetails.Mall),
		Address:       mallDetails.Address,
		Site:          mallDetails.Site,
		DayAndNight:   mallDetails.DayAndNight,
		WorkingHours:  workingHours,
		SubwayStation: subwayStation,
	}
	return mds
}
