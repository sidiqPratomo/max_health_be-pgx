package dto

import (
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/shopspring/decimal"
)

type UpdateDrugRequest struct {
	Name                   string          `json:"name" validate:"required"`
	GenericName            string          `json:"generic_name" validate:"required"`
	Content                string          `json:"content" validate:"required"`
	Manufacture            string          `json:"manufacture" validate:"required"`
	Description            string          `json:"description" validate:"required"`
	ClassificationId       int64           `json:"classification_id" validate:"required,number"`
	FormId                 int64           `json:"form_id" validate:"required,number"`
	CategoryId             int64           `json:"category_id" validate:"required,number"`
	UnitInPack             string          `json:"unit_in_pack" validate:"required"`
	SellingUnit            string          `json:"selling_unit" validate:"required"`
	Weight                 decimal.Decimal `json:"weight" validate:"required"`
	Height                 decimal.Decimal `json:"height" validate:"required"`
	Length                 decimal.Decimal `json:"length" validate:"required"`
	Width                  decimal.Decimal `json:"width" validate:"required"`
	IsActive               bool            `json:"is_active"`
	IsPrescriptionRequired bool            `json:"is_prescription_required"`
}

type CreateDrugRequest struct {
	Name                   string  `json:"name" binding:"required" validate:"required"`
	GenericName            string  `json:"generic_name" binding:"required" validate:"required"`
	Content                string  `json:"content" binding:"required" validate:"required"`
	Manufacture            string  `json:"manufacture" binding:"required" validate:"required"`
	Description            string  `json:"description" binding:"required" validate:"required"`
	ClassificationId       int64   `json:"classification_id" binding:"required" validate:"required"`
	FormId                 int64   `json:"form_id" binding:"required" validate:"required"`
	CategoryId             int64   `json:"category_id" binding:"required" validate:"required"`
	UnitInPack             string  `json:"unit_in_pack" binding:"required" validate:"required"`
	SellingUnit            string  `json:"selling_unit" binding:"required" validate:"required"`
	Weight                 float64 `json:"weight" binding:"required,gte=0" validate:"required,gte=0"`
	Height                 float64 `json:"height" binding:"required,gte=0" validate:"required,gte=0"`
	Length                 float64 `json:"length" binding:"required,gte=0" validate:"required,gte=0"`
	Width                  float64 `json:"width" binding:"required,gte=0" validate:"required,gte=0"`
	IsActive               bool    `json:"is_active"`
	IsPreScriptionRequired bool    `json:"is_prescription_required"`
}

type DrugDetailResponse struct {
	Id                     int64              `json:"id"`
	Name                   string             `json:"name"`
	GenericName            string             `json:"generic_name"`
	Content                string             `json:"content"`
	Manufacture            string             `json:"manufacture"`
	Description            string             `json:"description"`
	Classification         DrugClassification `json:"classification"`
	Form                   DrugForm           `json:"form"`
	Category               DrugCategory       `json:"category"`
	UnitInPack             string             `json:"unit_in_pack"`
	SellingUnit            string             `json:"selling_unit"`
	Weight                 decimal.Decimal    `json:"weight"`
	Height                 decimal.Decimal    `json:"height"`
	Length                 decimal.Decimal    `json:"length"`
	Width                  decimal.Decimal    `json:"width"`
	Image                  string             `json:"image"`
	IsActive               bool               `json:"is_active"`
	IsPrescriptionRequired bool               `json:"is_prescription_required"`
	PharmacyDrugs          []PharmacyDrug     `json:"pharmacy_drugs"`
}

type DrugListingResponse struct {
	Drugs    []entity.DrugListing `json:"drug_list"`
	PageInfo entity.PageInfo      `json:"page_info"`
}

type DrugResponse struct {
	Id                     int64              `json:"id"`
	Name                   string             `json:"name"`
	GenericName            string             `json:"generic_name"`
	Content                string             `json:"content"`
	Manufacture            string             `json:"manufacture"`
	Description            string             `json:"description"`
	Classification         DrugClassification `json:"classification"`
	Form                   DrugForm           `json:"form"`
	Category               DrugCategory       `json:"category"`
	UnitInPack             string             `json:"unit_in_pack"`
	SellingUnit            string             `json:"selling_unit"`
	Weight                 decimal.Decimal    `json:"weight"`
	Height                 decimal.Decimal    `json:"height"`
	Length                 decimal.Decimal    `json:"length"`
	Width                  decimal.Decimal    `json:"width"`
	Image                  string             `json:"image"`
	IsActive               bool               `json:"is_active"`
	IsPreScriptionRequired bool               `json:"is_prescription_required"`
}

type AllDrugsResponse struct {
	Drugs    []entity.Drug   `json:"drugs"`
	PageInfo entity.PageInfo `json:"page_info"`
}

type DrugClassification struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DrugForm struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DrugCategory struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"image"`
}

type PharmacyDrugByPharmacyDTO struct {
	Id    int64           `json:"pharmacy_drug_id"`
	Price decimal.Decimal `json:"price"`
	Stock int             `json:"stock"`
	Drug  DrugResponse    `json:"drug"`
}

type PharmacyDrugsByPharmacyResponse struct {
	Drugs    []PharmacyDrugByPharmacyDTO `json:"drugs"`
	PageInfo entity.PageInfo             `json:"page_info"`
}

type PrescriptionDrugItemRequest struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type PharmacyDrugMutationsResponse struct {
	Id              int64  `json:"pharmacy_drug_id"`
	PharmacyId      int64  `json:"pharmacy_id"`
	PharmacyName    string `json:"pharmacy_name"`
	PharmacyAddress string `json:"pharmacy_address"`
	Stock           int    `json:"stock"`
}

func ConvertToDrugClassificationListDTO(drugClassificationList []entity.DrugClassification) []DrugClassification {
	var drugClassificationListDTO []DrugClassification

	for _, drugClassification := range drugClassificationList {
		drugClassificationListDTO = append(drugClassificationListDTO, DrugClassification(drugClassification))
	}

	return drugClassificationListDTO
}

func ConvertToDrugFormListDTO(formList []entity.DrugForm) []DrugForm {
	var formListDTO []DrugForm

	for _, form := range formList {
		formListDTO = append(formListDTO, DrugForm(form))
	}

	return formListDTO
}

func ConvertToDrugDetailResponse(drug entity.Drug, pharmacyDrugList []entity.PharmacyDrug) DrugDetailResponse {
	return DrugDetailResponse{
		Id:                     drug.Id,
		Name:                   drug.Name,
		GenericName:            drug.GenericName,
		Content:                drug.Content,
		Manufacture:            drug.Manufacture,
		Description:            drug.Description,
		Classification:         DrugClassification(drug.Classification),
		Form:                   DrugForm(drug.Form),
		Category:               DrugCategory(drug.Category),
		UnitInPack:             drug.UnitInPack,
		SellingUnit:            drug.SellingUnit,
		Weight:                 drug.Weight,
		Height:                 drug.Height,
		Length:                 drug.Length,
		Width:                  drug.Width,
		Image:                  drug.Image,
		IsPrescriptionRequired: drug.IsPrescriptionRequired,
		PharmacyDrugs:          ConvertToPharmacyDrugListDTO(pharmacyDrugList),
	}
}

func ConvertToDrug(drugDTO UpdateDrugRequest) entity.Drug {
	return entity.Drug{
		Name:                   drugDTO.Name,
		GenericName:            drugDTO.GenericName,
		Content:                drugDTO.Content,
		Manufacture:            drugDTO.Manufacture,
		Description:            drugDTO.Description,
		Classification:         entity.DrugClassification{Id: drugDTO.ClassificationId},
		Form:                   entity.DrugForm{Id: drugDTO.FormId},
		Category:               entity.DrugCategory{Id: drugDTO.CategoryId},
		UnitInPack:             drugDTO.UnitInPack,
		SellingUnit:            drugDTO.SellingUnit,
		Weight:                 drugDTO.Weight,
		Height:                 drugDTO.Height,
		Length:                 drugDTO.Length,
		Width:                  drugDTO.Width,
		IsActive:               drugDTO.IsActive,
		IsPrescriptionRequired: drugDTO.IsPrescriptionRequired,
	}
}

func ConvertToDrugResponse(drug entity.Drug) DrugResponse {
	return DrugResponse{
		Id:                     drug.Id,
		Name:                   drug.Name,
		GenericName:            drug.GenericName,
		Content:                drug.Content,
		Manufacture:            drug.Manufacture,
		Description:            drug.Description,
		Classification:         DrugClassification(drug.Classification),
		Form:                   DrugForm(drug.Form),
		Category:               DrugCategory(drug.Category),
		UnitInPack:             drug.UnitInPack,
		SellingUnit:            drug.SellingUnit,
		Weight:                 drug.Weight,
		Height:                 drug.Height,
		Length:                 drug.Length,
		Width:                  drug.Width,
		Image:                  drug.Image,
		IsActive:               drug.IsActive,
		IsPreScriptionRequired: drug.IsPrescriptionRequired,
	}
}

func CreateDrugRequestToDrug(reqest CreateDrugRequest) entity.Drug {
	return entity.Drug{
		Name:                   reqest.Name,
		GenericName:            reqest.GenericName,
		Content:                reqest.Content,
		Manufacture:            reqest.Manufacture,
		Description:            reqest.Description,
		Classification:         entity.DrugClassification{Id: reqest.ClassificationId},
		Form:                   entity.DrugForm{Id: reqest.FormId},
		Category:               entity.DrugCategory{Id: reqest.CategoryId},
		UnitInPack:             reqest.UnitInPack,
		SellingUnit:            reqest.SellingUnit,
		Weight:                 decimal.NewFromFloat(reqest.Weight),
		Height:                 decimal.NewFromFloat(reqest.Height),
		Length:                 decimal.NewFromFloat(reqest.Length),
		Width:                  decimal.NewFromFloat(reqest.Width),
		IsActive:               reqest.IsActive,
		IsPrescriptionRequired: reqest.IsPreScriptionRequired,
	}
}

func ConvertPrescriptionDrugRequestToDrug(prescriptionDrugRequest PrescriptionDrugItemRequest) entity.Drug {
	return entity.Drug{
		Id:    prescriptionDrugRequest.Id,
		Name:  prescriptionDrugRequest.Name,
		Image: prescriptionDrugRequest.Image,
	}
}

func ConvertToMutationPharmacyDrugs(pharmacyDrugs []entity.PharmacyDrugDetail) []PharmacyDrugMutationsResponse {
	res := []PharmacyDrugMutationsResponse{}
	for _, pharmacyDrug := range pharmacyDrugs {
		res = append(res, PharmacyDrugMutationsResponse{Id: pharmacyDrug.Id, PharmacyName: pharmacyDrug.PharmacyName,
			PharmacyAddress: pharmacyDrug.PharmacyAddress,
			Stock:           pharmacyDrug.Stock, PharmacyId: pharmacyDrug.PharmacyId})
	}
	return res
}
