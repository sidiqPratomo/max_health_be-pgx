package dto

import (
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/shopspring/decimal"
)

type GetAllDoctorResponse struct {
	PageInfo entity.PageInfo `json:"page_info"`
	Doctors  []DoctorDto     `json:"doctors"`
}

type DoctorDto struct {
	DoctorId           int64           `json:"doctor_id"`
	AccountId          int64           `json:"account_id"`
	FeePerPatient      decimal.Decimal `json:"fee_per_patient"`
	IsOnline           bool            `json:"isOnline"`
	ProfilePicture     string          `json:"profile_picture"`
	Experience         int             `json:"experience"`
	Name               string          `json:"name" `
	SpecializationName string          `json:"specialization"`
}

type UpdateDoctorDataRequest struct {
	Name          string `json:"name" validate:"required"`
	Password      string `json:"password"`
	FeePerPatient string `json:"fee_per_patient" validate:"omitempty,number,gte=0"`
	Experience    int    `json:"years_of_experience" validate:"omitempty,number,gte=0"`
}

type DoctorSpecialization struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DoctorProfileResponse struct {
	Email              string          `json:"email"`
	Name               string          `json:"name"`
	ProfilePicture     string          `json:"profile_picture"`
	Experience         int             `json:"experience"`
	FeePerPatient      decimal.Decimal `json:"fee_per_patient"`
	SpecializationId   int64           `json:"specialization_id"`
	SpecializationName string          `json:"specialization_name"`
}

type UpdateDoctorStatusRequest struct {
	IsOnline bool `json:"is_online" binding:"required"`
}

type GetDoctorStatusResponse struct {
	IsOnline bool `json:"is_online"`
}

func UpdateDoctorDataRequestToDetailedDoctor(updateDoctorDataRequest UpdateDoctorDataRequest) entity.DetailedDoctor {
	decimal, _ := decimal.NewFromString(updateDoctorDataRequest.FeePerPatient)
	return entity.DetailedDoctor{
		Name:          updateDoctorDataRequest.Name,
		Password:      updateDoctorDataRequest.Password,
		FeePerPatient: decimal,
		Experience:    updateDoctorDataRequest.Experience,
	}
}

func ConvertToDocterSpecializationDTO(doctorSpecialization entity.DoctorSpecialization) DoctorSpecialization {
	return DoctorSpecialization{
		Id:   doctorSpecialization.Id,
		Name: doctorSpecialization.Name,
	}
}

func ConvertToDoctorSpecializationDTOList(doctorSpecializationList []entity.DoctorSpecialization) []DoctorSpecialization {
	var specializationListDTO []DoctorSpecialization

	for _, doctorSpecialization := range doctorSpecializationList {
		specializationListDTO = append(specializationListDTO, DoctorSpecialization{
			Id:   doctorSpecialization.Id,
			Name: doctorSpecialization.Name,
		})
	}

	return specializationListDTO
}
