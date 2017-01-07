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
type MallBase struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Logo       *Logo     `json:"logo"`
	Location   *Location `json:"location"`
	ShopsCount int       `json:"shops_count"`
}
type MallDetails struct {
	*MallBase
	Address       string         `json:"address"`
	Site          string         `json:"site"`
	DayAndNight   bool           `json:"day_and_night"`
	WorkingHours  []*WorkPeriod  `json:"working_hours"`
	SubwayStation *SubwayStation `json:"subway_staion"`
}

func serializeMallBase(mall *models.Mall) *MallBase {
	serializer := &MallBase{
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
		ShopsCount: mall.ShopsCount,
	}
	return serializer
}
func SerializeMall(mall *models.Mall) *MallDetails {
	workingHours := make([]*WorkPeriod, len(mall.WorkingHours))
	for i := range mall.WorkingHours {
		period := mall.WorkingHours[i]
		workingHours[i] = &WorkPeriod{
			Opening: &WeekTime{
				Time: period.OpenTime,
				Day:  period.OpenDay,
			},
			Closing: &WeekTime{
				Time: period.CloseTime,
				Day:  period.CloseDay,
			},
		}
	}
	var subwayStation *SubwayStation
	if mall.SubwayID != nil && mall.SubwayName != nil {
		subwayStation = &SubwayStation{ID: *mall.SubwayID, Name: *mall.SubwayName}
	}
	serializer := &MallDetails{
		MallBase:      serializeMallBase(mall),
		Address:       mall.Address,
		Site:          mall.Site,
		DayAndNight:   mall.DayAndNight,
		WorkingHours:  workingHours,
		SubwayStation: subwayStation,
	}
	return serializer
}
func SerializeMalls(malls []*models.Mall) []*MallBase {
	serializers := make([]*MallBase, len(malls))
	for i := range malls {
		mall := malls[i]
		serializers[i] = serializeMallBase(mall)
	}
	return serializers
}
