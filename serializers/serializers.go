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
	Phone       string `json:"phone"`
	Site        string `json:"site"`
	NearestMall *int   `json:"nearest_mall"`
}
type CategoryBase struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Logo       *Logo  `json:"logo"`
	ShopsCount int    `json:"shops_count"`
}
type CategoryDetails struct {
	*CategoryBase
}
type CityBase struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ShopsInMall struct {
	MallID  int   `json:"mall"`
	ShopIDs []int `json:"shops"`
}

type SearchResult struct {
	Mall     *MallBase `json:"mall"`
	ShopIDs  []int     `json:"shops"`
	Distance *float64  `json:"distance"`
}

func serializeMallBase(mall *models.Mall) *MallBase {
	if mall == nil {
		return nil
	}
	serializer := &MallBase{
		ID:    mall.ID,
		Name:  mall.Name,
		Phone: mall.Phone,
		Logo: &Logo{
			Large: mall.LogoLarge,
			Small: mall.LogoSmall,
		},
		Location: &Location{
			Lat: mall.Location.Lat,
			Lon: mall.Location.Lon,
		},
		ShopsCount: mall.ShopsCount,
	}
	return serializer
}
func serializeShopBase(shop *models.Shop) *ShopBase {
	if shop == nil {
		return nil
	}
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
func serializeCategoryBase(category *models.Category) *CategoryBase {
	if category == nil {
		return nil
	}
	serializer := &CategoryBase{
		ID:   category.ID,
		Name: category.Name,
		Logo: &Logo{
			Large: category.LogoLarge,
			Small: category.LogoSmall,
		},
		ShopsCount: category.ShopsCount,
	}
	return serializer
}
func serializeCityBase(city *models.City) *CityBase {
	if city == nil {
		return nil
	}
	serializer := &CityBase{
		ID:   city.ID,
		Name: city.Name,
	}
	return serializer
}
func SerializeMall(mall *models.Mall) *MallDetails {
	if mall == nil {
		return nil
	}
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
	if shop == nil {
		return nil
	}
	serializer := &ShopDetails{
		ShopBase:    serializeShopBase(shop),
		Phone:       shop.Phone,
		Site:        shop.Site,
		NearestMall: shop.NearestMall,
	}
	return serializer
}
func SerializeCategory(category *models.Category) *CategoryDetails {
	if category == nil {
		return nil
	}
	serializer := &CategoryDetails{
		CategoryBase: serializeCategoryBase(category),
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
func SerializeCategories(categories []*models.Category) []*CategoryBase {
	serializers := make([]*CategoryBase, len(categories))
	for i := range categories {
		category := categories[i]
		serializers[i] = serializeCategoryBase(category)
	}
	return serializers
}
func SerializeCities(cities []*models.City) []*CityBase {
	serializers := make([]*CityBase, len(cities))
	for i := range cities {
		city := cities[i]
		serializers[i] = serializeCityBase(city)
	}
	return serializers
}
func SerializeShopsInMalls(mallsShops models.ShopsInMalls) []*ShopsInMall {
	serializers := []*ShopsInMall{}
	for mallID, shopIDs := range mallsShops {
		if shopIDs == nil {
			shopIDs = []int{}
		}
		serializers = append(serializers, &ShopsInMall{MallID: mallID, ShopIDs: shopIDs})
	}
	return serializers
}
func SerializeSearchResults(searchResults []*models.SearchResult) []*SearchResult {
	serializers := make([]*SearchResult, len(searchResults))
	for i := range searchResults {
		searchResult := searchResults[i]
		shopIDs := searchResult.ShopIDs
		if shopIDs == nil {
			shopIDs = []int{}
		}
		serializers[i] = &SearchResult{
			Mall:     serializeMallBase(searchResult.Mall),
			ShopIDs:  shopIDs,
			Distance: searchResult.Distance,
		}
	}
	return serializers
}
