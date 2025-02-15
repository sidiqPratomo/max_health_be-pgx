package dto

import (
	"strings"

	"github.com/shopspring/decimal"

	"github.com/sidiqPratomo/max-health-backend/entity"
)

type UserAddressResponse struct {
	Id          int64               `json:"id"`
	Province    ProvinceResponse    `json:"province"`
	City        CityResponse        `json:"city"`
	District    DistrictResponse    `json:"district"`
	Subdistrict SubdistrictResponse `json:"subdistrict"`
	Latitude    string              `json:"latitude"`
	Longitude   string              `json:"longitude"`
	Label       string              `json:"label"`
	Address     string              `json:"address"`
	IsMain      bool                `json:"is_main"`
}

type UserProfileResponse struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	GenderId       int64  `json:"gender_id"`
	Gender         string `json:"gender"`
	DateOfBirth    string `json:"date_of_birth"`
}

type UpdateUserDataRequest struct {
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password"`
	GenderId    int64  `json:"gender_id" validate:"required,number,gte=1"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

type AllUserAddressResponse struct {
	Address []UserAddressResponse `json:"address"`
}

type AddUserAddressRequest struct {
	ProvinceId    int64  `json:"province_id" binding:"required,gte=1,lte=37"`
	CityId        int64  `json:"city_id" binding:"required,gte=1,lte=514"`
	DistrictId    int64  `json:"district_id" binding:"required,gte=1,lte=7277"`
	SubdistrictId int64  `json:"subdistrict_id" binding:"required,gte=1,lte=83761"`
	Latitude      string `json:"latitude" binding:"required,latitude"`
	Longitude     string `json:"longitude" binding:"required,longitude"`
	Label         string `json:"label"`
	Address       string `json:"address" binding:"required"`
	IsMain        bool   `json:"is_main"`
}

type AddUserAddressAutofillRequest struct {
	ProvinceName    string `json:"province_name" binding:"required"`
	CityName        string `json:"city_name" binding:"required"`
	DistrictName    string `json:"district_name" binding:"required"`
	SubdistrictName string `json:"subdistrict_name" binding:"required"`
	Latitude        string `json:"latitude" binding:"required,latitude"`
	Longitude       string `json:"longitude" binding:"required,longitude"`
	Label           string `json:"label"`
	Address         string `json:"address" binding:"required"`
	IsMain          bool   `json:"is_main"`
}

type UpdateUserAddressRequest struct {
	ProvinceId    int64  `json:"province_id" binding:"required,gte=1,lte=37"`
	CityId        int64  `json:"city_id" binding:"required,gte=1,lte=514"`
	DistrictId    int64  `json:"district_id" binding:"required,gte=1,lte=7277"`
	SubdistrictId int64  `json:"subdistrict_id" binding:"required,gte=1,lte=83761"`
	Latitude      string `json:"latitude" binding:"required,latitude"`
	Longitude     string `json:"longitude" binding:"required,longitude"`
	Label         string `json:"label"`
	Address       string `json:"address" binding:"required"`
	IsActive      bool   `json:"is_active"`
	IsMain        bool   `json:"is_main"`
}

type CreateOneCartRequest struct {
	PharmacyDrugId int64 `json:"pharmacy_drug_id" binding:"required"`
}

type UpdateQtyCartRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

type CartDTOResponse struct {
	Page  entity.PageInfo `json:"page_info"`
	Carts []CartDTO       `json:"carts"`
}

type CartDTO struct {
	Id            int64            `json:"cart_item_id"`
	UserId        int64            `json:"user_id"`
	Quantity      int              `json:"quantity"`
	PharmacyDrugs PharmacyDrugsDto `json:"pharmacy_drugs"`
}

type PharmacyDrugsDto struct {
	Id           int64           `json:"pharmacy_drug_id"`
	PharmacyId   int64           `json:"pharmacy_id"`
	Name         string          `json:"drug_name"`
	Price        decimal.Decimal `json:"price"`
	Image        string          `json:"image"`
	PharmacyName string          `json:"pharmacy_name"`
	Stock        int             `json:"stock"`
}

func ConvertUpdateRequestToUserAddress(request UpdateUserAddressRequest) entity.UserAddress {
	return entity.UserAddress{
		Province:    entity.Province{Id: request.ProvinceId},
		City:        entity.City{Id: request.CityId},
		District:    entity.District{Id: request.DistrictId},
		Subdistrict: entity.Subdistrict{Id: request.SubdistrictId},
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		Label:       request.Label,
		Address:     request.Address,
		IsActive:    request.IsActive,
		IsMain:      request.IsMain,
	}
}

func ConvertAddRequestToUserAddress(request AddUserAddressRequest) entity.UserAddress {
	return entity.UserAddress{
		Province:    entity.Province{Id: request.ProvinceId},
		City:        entity.City{Id: request.CityId},
		District:    entity.District{Id: request.DistrictId},
		Subdistrict: entity.Subdistrict{Id: request.SubdistrictId},
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		Label:       request.Label,
		Address:     request.Address,
		IsMain:      request.IsMain,
	}
}

func ConvertUserAddressToAddUserAddressRequest(userAddress entity.UserAddress) AddUserAddressRequest {
	return AddUserAddressRequest{
		ProvinceId:    userAddress.Province.Id,
		CityId:        userAddress.City.Id,
		DistrictId:    userAddress.District.Id,
		SubdistrictId: userAddress.Subdistrict.Id,
		Latitude:      userAddress.Latitude,
		Longitude:     userAddress.Longitude,
		Label:         userAddress.Label,
		Address:       userAddress.Address,
		IsMain:        userAddress.IsMain,
	}
}

func ConvertAddRequestAutofillToUserAddress(request AddUserAddressAutofillRequest) entity.UserAddress {
	return entity.UserAddress{
		ProvinceName:    strings.Trim(request.ProvinceName, ""),
		CityName:        strings.Trim(request.CityName, ""),
		DistrictName:    strings.Trim(request.DistrictName, ""),
		SubdistrictName: strings.Trim(request.SubdistrictName, ""),
		Latitude:        request.Latitude,
		Longitude:       request.Longitude,
		Label:           request.Label,
		Address:         request.Address,
		IsMain:          request.IsMain,
	}
}

func UserAddressToUserAddressResponse(userAddress entity.UserAddress) UserAddressResponse {
	return UserAddressResponse{
		Id:          userAddress.Id,
		Province:    ProvinceToProvinceResponse(userAddress.Province),
		City:        CityToCityResponse(userAddress.City),
		District:    DistrictToDistrictResponse(userAddress.District),
		Subdistrict: SubdistrictToSubdistrictResponse(userAddress.Subdistrict),
		Latitude:    userAddress.Latitude,
		Longitude:   userAddress.Longitude,
		Label:       userAddress.Label,
		Address:     userAddress.Address,
		IsMain:      userAddress.IsMain,
	}
}

func ConvertToAllUserAddressResponse(userAddress []entity.UserAddress) *AllUserAddressResponse {
	data := []UserAddressResponse{}
	for _, address := range userAddress {
		data = append(data, UserAddressToUserAddressResponse(address))
	}
	return &AllUserAddressResponse{
		Address: data,
	}
}

func UpdateUserDataRequestToDetailedUser(updateUserDataRequest UpdateUserDataRequest) entity.DetailedUser {
	return entity.DetailedUser{
		Name:        updateUserDataRequest.Name,
		Password:    updateUserDataRequest.Password,
		GenderId:    updateUserDataRequest.GenderId,
		DateOfBirth: updateUserDataRequest.DateOfBirth,
	}
}
