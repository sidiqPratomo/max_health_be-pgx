package dto

import (
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/shopspring/decimal"
)

type PharmacyDrug struct {
	Id               int64           `json:"id"`
	Pharmacy         Pharmacy        `json:"pharmacy"`
	DrugId           int64           `json:"drug_id"`
	Price            decimal.Decimal `json:"price"`
	Stock            int             `json:"stock"`
	CartItemId       *int            `json:"cart_item_id,omitempty"`
	CartItemQuantity *int            `json:"cart_item_quantity,omitempty"`
}

type PharmacyRequest struct {
	Name              string `json:"pharmacy_name" binding:"required"`
	PharmacyManagerId int64  `json:"pharmacy_manager_id" binding:"required"`
}

type UpdatePharmacyRequest struct {
	Id                      int64                              `json:"id"`
	Name                    string                             `json:"pharmacy_name" binding:"required"`
	PharmacistName          string                             `json:"pharmacist_name" binding:"required"`
	PharmacistLicenseNumber string                             `json:"pharmacist_license_number" binding:"required"`
	PharmacistPhoneNumber   string                             `json:"pharmacist_phone_number" binding:"required"`
	City                    string                             `json:"city" binding:"required"`
	Address                 string                             `json:"address" binding:"required"`
	Latitude                string                             `json:"latitude" binding:"required,latitude"`
	Longitude               string                             `json:"longitude" binding:"required,longitude"`
	Operationals            []UpdatePharmacyOperationalRequest `json:"operationals" binding:"required"`
	Couriers                []UpdatePharmacyCourierRequest     `json:"couriers" binding:"required"`
}

type UpdatePharmacyOperationalRequest struct {
	Id             int64  `json:"id" binding:"required"`
	OperationalDay string `json:"operational_day" binding:"required"`
	OpenHour       string `json:"open_hour" binding:"required"`
	CloseHour      string `json:"close_hour" binding:"required"`
	IsOpen         bool   `json:"is_open"`
}

type UpdatePharmacyCourierRequest struct {
	Id        int64 `json:"id" binding:"required"`
	CourierId int64 `json:"courier_id" binding:"required"`
	IsActive  bool  `json:"is_active"`
}

type Pharmacy struct {
	Id                      int64   `json:"id,omitempty"`
	PharmacyManagerId       int64   `json:"manager_id,omitempty"`
	Name                    string  `json:"pharmacy_name,omitempty"`
	PharmacistName          string  `json:"pharmacist_name,omitempty"`
	PharmacistLicenseNumber string  `json:"pharamcist_license_name,omitempty"`
	PharmacistPhoneNumber   string  `json:"pharmacist_phone_number,omitempty"`
	City                    string  `json:"city,omitempty"`
	Address                 string  `json:"address,omitempty"`
	Latitude                string  `json:"latitude,omitempty"`
	Longitude               string  `json:"longitude,omitempty"`
	Distance                float64 `json:"distance,omitempty"`
}

type UpdatePharmacyDrugReq struct {
	Stock int             `json:"stock"`
	Price decimal.Decimal `json:"price"`
}

type AddPharmacyDrugReq struct {
	PharmacyId int64           `json:"pharmacy_id"`
	DrugId     int64           `json:"drug_id"`
	Stock      int             `json:"stock"`
	Price      decimal.Decimal `json:"price"`
}

type PostStockMutationRequest struct {
	RecipientPharmacyDrugId int64
	SenderPharmacyDrugId    int64 `json:"sender_pharmacy_drug_id" binding:"required"`
	Quantity                int   `json:"quantity" binding:"required,gte=1"`
}

type GetAllPharmacyResponse struct {
	PageInfo   entity.PageInfo `json:"page_info"`
	Pharmacies []PharmacyDTO   `json:"pharmacies"`
}

type PharmacyDTO struct {
	Id                      int64  `json:"id,omitempty"`
	PharmacyManagerId       int64  `json:"manager_id,omitempty"`
	Name                    string `json:"pharmacy_name,omitempty"`
	PharmacistName          string `json:"pharmacist_name,omitempty"`
	PharmacistLicenseNumber string `json:"pharamcist_license_name,omitempty"`
	PharmacistPhoneNumber   string `json:"pharmacist_phone_number,omitempty"`
	City                    string `json:"city,omitempty"`
	Address                 string `json:"address,omitempty"`
	Latitude                string `json:"latitude,omitempty"`
	Longitude               string `json:"longitude,omitempty"`
}

type DetailPharmacyDrug struct {
	Id    int64           `json:"id"`
	Drug  DrugResponse    `json:"drug"`
	Price decimal.Decimal `json:"price"`
}

func ConvertToPharmacyDrugDTO(pharmacyDrug entity.PharmacyDrug) PharmacyDrug {
	return PharmacyDrug{
		Id: pharmacyDrug.Id,
		Pharmacy: Pharmacy{
			Id:                      pharmacyDrug.Pharmacy.Id,
			Name:                    pharmacyDrug.Pharmacy.Name,
			PharmacistName:          pharmacyDrug.Pharmacy.PharmacistName,
			PharmacistLicenseNumber: pharmacyDrug.Pharmacy.PharmacistLicenseNumber,
			PharmacistPhoneNumber:   pharmacyDrug.Pharmacy.PharmacistPhoneNumber,
			Distance:                pharmacyDrug.Pharmacy.Distance},
		DrugId:           pharmacyDrug.DrugId,
		Price:            pharmacyDrug.Price,
		Stock:            pharmacyDrug.Stock,
		CartItemId:       pharmacyDrug.CartItemId,
		CartItemQuantity: pharmacyDrug.CartItemQuantity,
	}
}

func ConvertToPharmacyDrugListDTO(pharmacyDrugList []entity.PharmacyDrug) []PharmacyDrug {
	var pharmacyDrugListDTO []PharmacyDrug

	for _, pharmacyDrug := range pharmacyDrugList {
		pharmacyDrugDTO := ConvertToPharmacyDrugDTO(pharmacyDrug)

		pharmacyDrugListDTO = append(pharmacyDrugListDTO, pharmacyDrugDTO)
	}

	return pharmacyDrugListDTO
}

func ConvertPharmacyRequestToPharmacy(request PharmacyRequest) entity.Pharmacy {
	return entity.Pharmacy{
		PharmacyManagerId: request.PharmacyManagerId,
		Name:              request.Name,
	}
}

func ConvertUpdatePharmacyRequestToPharmacy(request UpdatePharmacyRequest) entity.Pharmacy {
	return entity.Pharmacy{
		Id:                      request.Id,
		Name:                    request.Name,
		PharmacistName:          request.PharmacistName,
		PharmacistLicenseNumber: request.PharmacistLicenseNumber,
		PharmacistPhoneNumber:   request.PharmacistPhoneNumber,
		City:                    request.City,
		Latitude:                request.Latitude,
		Longitude:               request.Longitude,
		Address:                 request.Address,
	}
}

func UpdatePharmacyOperationalRequestToPharmacyOperational(request UpdatePharmacyOperationalRequest) entity.PharmacyOperational {
	return entity.PharmacyOperational{
		Id:             request.Id,
		OperationalDay: request.OperationalDay,
		OpenHour:       request.OpenHour,
		CloseHour:      request.CloseHour,
		IsOpen:         request.IsOpen,
	}
}

func AllUpdatePharmacyOperationalsToAllPharmacyOperationals(updatePharmacyOperationalRequests []UpdatePharmacyOperationalRequest) []entity.PharmacyOperational {
	var pharmacyOperationals []entity.PharmacyOperational

	for _, updatePharmacyOperationalRequest := range updatePharmacyOperationalRequests {
		pharmacyOperationals = append(pharmacyOperationals, UpdatePharmacyOperationalRequestToPharmacyOperational(updatePharmacyOperationalRequest))
	}

	return pharmacyOperationals
}

func UpdatePharmacyCourierRequestToPharmacyCourier(request UpdatePharmacyCourierRequest) entity.PharmacyCourier {
	return entity.PharmacyCourier{
		Id:        request.Id,
		CourierId: request.CourierId,
		IsActive:  request.IsActive,
	}
}

func AllUpdatePharmacyCourierRequestToAllPharmacyCouriers(pharmacyCourierRequests []UpdatePharmacyCourierRequest) []entity.PharmacyCourier {
	var pharmacyCouriers []entity.PharmacyCourier

	for _, pharmacyCourierRequest := range pharmacyCourierRequests {
		pharmacyCouriers = append(pharmacyCouriers, UpdatePharmacyCourierRequestToPharmacyCourier(pharmacyCourierRequest))
	}

	return pharmacyCouriers
}
