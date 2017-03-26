package models

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
	MallID  int   `json:"mall"`
	ShopIDs []int `json:"shops"`
}

type SearchResult struct {
	Mall     *Mall
	ShopIDs  []int
	Distance *float64
}
