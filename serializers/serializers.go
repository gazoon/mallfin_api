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
type ShopBase struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Logo       *Logo  `json:"logo"`
	Score      int    `json:"score"`
	MallsCount int    `json:"malls_count"`
}
type ShopDetails struct {
	*ShopBase
	Phone string `json:"phone"`
	Site  string `json:"site"`
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
func serializeShopBase(shop *models.Shop) *ShopBase {
	serializer := &ShopBase{
		ID:   shop.ID,
		Name: shop.Name,
		Logo: &Logo{
			Large: shop.LogoLarge,
			Small: shop.LogoSmall,
		},
		Score:      shop.Score,
		MallsCount: shop.MallsCount,
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
func SerializeShop(shop *models.Shop) *ShopDetails {
	serializer := &ShopDetails{
		ShopBase: serializeShopBase(shop),
		Phone:    shop.Phone,
		Site:     shop.Site,
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
func SerializeShops(shops []*models.Shop) []*ShopBase {
	serializers := make([]*ShopBase, len(shops))
	for i := range shops {
		shop := shops[i]
		serializers[i] = serializeShopBase(shop)
	}
	return serializers
}
