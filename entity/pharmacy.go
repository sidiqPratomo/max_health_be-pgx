package entity

import (
	"github.com/shopspring/decimal"
)

type PharmacyManager struct {
	Id      int64
	Account Account
}

type PharmacyDrug struct {
	Id               int64
	Pharmacy         Pharmacy
	DrugId           int64
	Price            decimal.Decimal
	Stock            int
	CartItemId       *int
	CartItemQuantity *int
}

type DrugListing struct {
	PharmacyDrugId         int64           `json:"pharmacy_drug_id"`
	DrugId                 int64           `json:"drug_id"`
	Name                   string          `json:"drug_name"`
	MinPrice               decimal.Decimal `json:"min_price"`
	MaxPrice               decimal.Decimal `json:"max_price"`
	Image                  string          `json:"image_url"`
	IsPrescriptionRequired string          `json:"prescription_required"`
}

type PharmacyOperational struct {
	Id             int64
	PharmacyId     int64
	OperationalDay string
	OpenHour       string
	CloseHour      string
	IsOpen         bool
}

type Pharmacy struct {
	Id                      int64
	PharmacyManagerId       int64
	Name                    string
	PharmacistName          string
	PharmacistLicenseNumber string
	PharmacistPhoneNumber   string
	City                    string
	Address                 string
	Latitude                string
	Longitude               string
	Distance                float64
}

type PharmacyJoinPharmacyDrug struct {
	Id                      int64
	PharmacyManagerId       int64
	Name                    string
	PharmacistName          string
	PharmacistLicenseNumber string
	PharmacistPhoneNumber   string
	City                    string
	Address                 string
	Latitude                string
	Longitude               string
	PharmacyDrug            []PharmacyDrug
}

type PharmacyCourier struct {
	Id         int64
	PharmacyId int64
	CourierId  int64
	IsActive   bool
}

type Courier struct {
	Id         int64
	Name       string
	Price      decimal.Decimal
	IsOfficial bool
}

type PharmacyDrugDetail struct {
	Id              int64
	PharmacyId      int64
	PharmacyName    string
	PharmacyAddress string
	DrugId          int64
	Price           decimal.Decimal
	Stock           int
}

type CourierOption struct {
	Price float64 `json:"price"`
	Etd   string  `json:"estimated_time_of_delivery"`
}

type AvailableCourier struct {
	PharmacyCourierId int64           `json:"pharmacy_courier_id"`
	CourierName       string          `json:"courier_name"`
	CourierOptions    []CourierOption `json:"options"`
}

type PharmacyDeliveryFee struct {
	Id           int64              `json:"pharmacy_id"`
	PharmacyName string             `json:"pharmacy_name"`
	Distance     int                `json:"distance"`
	Couriers     []AvailableCourier `json:"couriers"`
}

type PharmacyDrugByPharmacyId struct {
	Id    int64           `json:"pharmacy_drug_id"`
	Price decimal.Decimal `json:"price"`
	Stock int             `json:"stock"`
	Drug  Drug            `json:"drug"`
}

type PharmacyDrugAndCartId struct {
	CartId         int64
	PharmacyDrugId int64
	Stock          int
}

type DetailPharmacyDrug struct {
	Id    int64
	Drug  Drug
	Price decimal.Decimal
}
