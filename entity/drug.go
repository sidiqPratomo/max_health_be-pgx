package entity

import (
	"github.com/shopspring/decimal"
)

type Drug struct {
	Id                     int64
	Name                   string
	GenericName            string
	Content                string
	Manufacture            string
	Description            string
	Classification         DrugClassification
	Form                   DrugForm
	Category               DrugCategory
	UnitInPack             string
	SellingUnit            string
	Weight                 decimal.Decimal
	Height                 decimal.Decimal
	Length                 decimal.Decimal
	Width                  decimal.Decimal
	Image                  string
	IsActive               bool
	IsPrescriptionRequired bool
}

type DrugDetail struct {
	Id                     int64
	Name                   string
	GenericName            string
	Content                string
	Manufacture            string
	Description            string
	ClassificationId       int64
	FormId                 int64
	UnitInPack             string
	SellingUnit            string
	Weight                 decimal.Decimal
	Height                 decimal.Decimal
	Length                 decimal.Decimal
	Width                  decimal.Decimal
	Image                  string
	DrugCategoryId         int64
	IsPrescriptionRequired bool
	IsActive               bool
}

type DrugClassification struct {
	Id   int64
	Name string
}

type DrugForm struct {
	Id   int64
	Name string
}

type DrugCategory struct {
	Id   int64
	Name string
	Url  string
}

type DrugQuantity struct {
	PharmacyDrug DetailPharmacyDrug
	Quantity     int
}

type PrepareForCheckoutItem struct {
	PharmacyId      int64
	PharmacyName    string
	PharmacyAddress string
	Distance        float64
	Subtotal        decimal.Decimal
	DeliveryOptions []AvailableCourier
	Weight          decimal.Decimal
	DrugQuantities  []DrugQuantity
}

type PrepareForCheckout struct {
	Items       []PrepareForCheckoutItem
	UserAddress UserAddress
}
