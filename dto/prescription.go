package dto

import (
	"time"

	"github.com/sidiqPratomo/max-health-backend/entity"
)

type PrescriptionDrugRequest struct {
	Id       int64
	Drug     PrescriptionDrugItemRequest `json:"drug"`
	Quantity int                         `json:"quantity" validate:"required,gte=1"`
	Note     string                      `json:"note"`
}

type PrescriptionDrugResponse struct {
	Id         int64        `json:"id"`
	Drug       DrugResponse `json:"drug"`
	Quantity   int          `json:"quantity"`
	Note       string       `json:"note"`
	RedeemedAt *string      `json:"redeemed_at,omitempty"`
	OrderedAt  *string      `json:"ordered_at,omitempty"`
}

type PrescriptionResponse struct {
	Id                *int64                     `json:"id"`
	UserAccountId     int64                      `json:"user_account_id"`
	UserName          string                     `json:"user_name"`
	DoctorAccountId   int64                      `json:"doctor_account_id"`
	DoctorName        string                     `json:"doctor_name"`
	RedeemedAt        *time.Time                 `json:"redeemed_at"`
	OrderedAt         *time.Time                 `json:"ordered_at"`
	CreatedAt         *time.Time                 `json:"created_at"`
	PrescriptionDrugs []PrescriptionDrugResponse `json:"prescription_drugs"`
}

type PrescriptionResponseList struct {
	Prescriptions []PrescriptionResponse `json:"prescriptions"`
	PageInfo      struct {
		TotalPage int `json:"total_page"`
		TotalItem int `json:"total_item"`
	} `json:"page_info"`
}

func ConvertToPrescriptionDrugListResponse(prescriptionDrugList []entity.PrescriptionDrug) []PrescriptionDrugResponse {
	var list []PrescriptionDrugResponse

	for _, prescriptionDrug := range prescriptionDrugList {
		list = append(list, ConvertToPrescriptionDrugResponse(prescriptionDrug))
	}

	return list
}

func ConvertToPrescriptionDrugResponse(prescriptionDrug entity.PrescriptionDrug) PrescriptionDrugResponse {
	return PrescriptionDrugResponse{
		Id:       prescriptionDrug.Id,
		Drug:     ConvertToDrugResponse(prescriptionDrug.Drug),
		Quantity: prescriptionDrug.Quantity,
		Note:     prescriptionDrug.Note,
	}
}

func ConvertToPrescriptionResponse(prescription entity.Prescription) PrescriptionResponse {
	return PrescriptionResponse{
		Id:                prescription.Id,
		UserAccountId:     prescription.UserAccountId,
		UserName:          prescription.UserName,
		DoctorAccountId:   prescription.DoctorAccountId,
		DoctorName:        prescription.DoctorName,
		RedeemedAt:        prescription.RedeemedAt,
		OrderedAt:         prescription.OrderedAt,
		CreatedAt:         prescription.CreatedAt,
		PrescriptionDrugs: ConvertToPrescriptionDrugListResponse(prescription.PrescriptionDrugs),
	}
}

func ConvertToPrescriptionResponseList(prescriptionList []entity.Prescription, totalItem, totalPage int) PrescriptionResponseList {
	var responseList []PrescriptionResponse

	for _, prescription := range prescriptionList {
		responseList = append(responseList, ConvertToPrescriptionResponse(prescription))
	}

	return PrescriptionResponseList{
		Prescriptions: responseList,
		PageInfo: struct {
			TotalPage int "json:\"total_page\""
			TotalItem int "json:\"total_item\""
		}{TotalPage: totalPage, TotalItem: totalItem},
	}
}
