package entity

type Province struct {
	Id   int64
	Code string
	Name string
}

type City struct {
	Id           int64
	Code         string
	ProvinceCode string
	Name         string
}

type District struct {
	Id       int64
	Code     string
	CityCode string
	Name     string
}

type Subdistrict struct {
	Id           int64
	Code         string
	DistrictCode string
	Name         string
}
