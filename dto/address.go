package dto

import "github.com/sidiqPratomo/max-health-backend/entity"

type AddressQuery struct {
	ProvinceCode string `form:"province_code"`
	CityCode     string `form:"city_code"`
	DistrictCode string `form:"district_code"`
}

type ProvinceResponse struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type AllProvincesResponse struct {
	Provinces []ProvinceResponse `json:"provinces"`
}

func ProvinceToProvinceResponse(province entity.Province) ProvinceResponse {
	return ProvinceResponse{
		Id:   province.Id,
		Code: province.Code,
		Name: province.Name,
	}
}

func ConvertToAllProvincesResponse(provinces []entity.Province) *AllProvincesResponse {
	data := []ProvinceResponse{}
	for _, province := range provinces {
		data = append(data, ProvinceToProvinceResponse(province))
	}
	return &AllProvincesResponse{
		Provinces: data,
	}
}

type CityResponse struct {
	Id           int64  `json:"id"`
	Code         string `json:"code"`
	ProvinceCode string `json:"province_code"`
	Name         string `json:"name"`
}

type AllCitiesResponse struct {
	Cities []CityResponse `json:"cities"`
}

func CityToCityResponse(city entity.City) CityResponse {
	return CityResponse{
		Id:           city.Id,
		Code:         city.Code,
		ProvinceCode: city.ProvinceCode,
		Name:         city.Name,
	}
}

func ConvertToAllCitiesResponse(cities []entity.City) *AllCitiesResponse {
	data := []CityResponse{}
	for _, city := range cities {
		data = append(data, CityToCityResponse(city))
	}
	return &AllCitiesResponse{
		Cities: data,
	}
}

type DistrictResponse struct {
	Id       int64  `json:"id"`
	Code     string `json:"code"`
	CityCode string `json:"city_code"`
	Name     string `json:"name"`
}

type AllDistrictsResponse struct {
	Districts []DistrictResponse `json:"districts"`
}

func DistrictToDistrictResponse(district entity.District) DistrictResponse {
	return DistrictResponse{
		Id:       district.Id,
		Code:     district.Code,
		CityCode: district.CityCode,
		Name:     district.Name,
	}
}

func ConvertToAllDistrictsResponse(districts []entity.District) *AllDistrictsResponse {
	data := []DistrictResponse{}
	for _, district := range districts {
		data = append(data, DistrictToDistrictResponse(district))
	}
	return &AllDistrictsResponse{
		Districts: data,
	}
}

type SubdistrictResponse struct {
	Id           int64  `json:"id"`
	Code         string `json:"code"`
	DistrictCode string `json:"district_code"`
	Name         string `json:"name"`
}

type AllSubdistrictsResponse struct {
	Subdistricts []SubdistrictResponse `json:"subdistricts"`
}

func SubdistrictToSubdistrictResponse(subdistrict entity.Subdistrict) SubdistrictResponse {
	return SubdistrictResponse{
		Id:           subdistrict.Id,
		Code:         subdistrict.Code,
		DistrictCode: subdistrict.DistrictCode,
		Name:         subdistrict.Name,
	}
}

func ConvertToAllSubdistrictsResponse(subdistricts []entity.Subdistrict) *AllSubdistrictsResponse {
	data := []SubdistrictResponse{}
	for _, subdistrict := range subdistricts {
		data = append(data, SubdistrictToSubdistrictResponse(subdistrict))
	}
	return &AllSubdistrictsResponse{
		Subdistricts: data,
	}
}
