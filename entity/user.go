package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	Id          int64
	AccountId   int64
	GenderId    *int64
	GenderName  *string
	DateOfBirth *time.Time
}

type DetailedUser struct {
	Id             int64
	Name           string
	ProfilePicture string
	Password       string
	GenderId       int64
	DateOfBirth    string
}

type UserAddress struct {
	Id              int64
	UserId          int64
	Province        Province
	City            City
	District        District
	Subdistrict     Subdistrict
	ProvinceName    string
	CityName        string
	DistrictName    string
	SubdistrictName string
	Latitude        string
	Longitude       string
	Label           string
	Address         string
	IsActive        bool
	IsMain          bool
}

type CartItem struct {
	Id             int64
	UserId         int64
	PharmacyDrugId int64
	Quantity       int
}

type CartItemForCheckout struct {
	Id             int64
	DrugId         int64
	DrugName       string
	PharmacyDrugId int64
	Price          int
	Unit           string
	Quantity       int
}

type CartItemChanges struct {
	PharmacyDrugId int64
	Quantity       int
	Stock          int
}

type CartItemData struct {
	CartItemId     int64
	UserId         int64
	PharmacyDrugId int64
	Quantity       int
	PharmacyId     int64
	DrugId         int64
	Price          decimal.Decimal
	Stock          int
	DrugName       string
	Image          string
	PharmacyName   string
}
