package dto

import "github.com/sidiqPratomo/max-health-backend/entity"

type DeliveryFeeRequest struct {
	UserAddressId int64   `json:"user_address_id" binding:"required,gte=1"`
	CartItemsId   []int64 `json:"cart_items_id" binding:"required"`
	AccountId     int64
}

type AllDeliveryFeeResponse struct {
	Pharmacies []entity.PharmacyDeliveryFee `json:"pharmacies"`
}
