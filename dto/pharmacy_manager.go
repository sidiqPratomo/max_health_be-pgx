package dto

import "github.com/sidiqPratomo/max-health-backend/entity"

type PartnerResponse struct {
	Id             int64  `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
}

type AllPartnersResponse struct {
	Partners []PartnerResponse `json:"partners"`
}

func PharmacyManagerToPartnerResponse(pharmacyManager entity.PharmacyManager) PartnerResponse {
	return PartnerResponse{
		Id:             pharmacyManager.Id,
		Email:          pharmacyManager.Account.Email,
		Name:           pharmacyManager.Account.Name,
		ProfilePicture: pharmacyManager.Account.ProfilePicture,
	}
}

func ConvertToAllPartnersResponse(pharmacyManagers []entity.PharmacyManager) *AllPartnersResponse {
	data := []PartnerResponse{}
	for _, pharmacyManager := range pharmacyManagers {
		data = append(data, PharmacyManagerToPartnerResponse(pharmacyManager))
	}
	return &AllPartnersResponse{
		Partners: data,
	}
}
