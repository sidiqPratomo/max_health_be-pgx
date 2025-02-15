package dto

import (
	"time"

	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/shopspring/decimal"
)

type AllOrdersResponse struct {
	PageInfo entity.PageInfo `json:"page_info"`
	Orders   []OrderResponse `json:"orders"`
}

type OrderResponse struct {
	Id             int64                   `json:"id"`
	Address        string                  `json:"address,omitempty"`
	PaymentProof   string                  `json:"payment_proof,omitempty"`
	TotalAmount    decimal.Decimal         `json:"total_amount"`
	PharmacyOrders []OrderPharmacyResponse `json:"pharmacies"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

type AllOrderPharmaciesResponse struct {
	PageInfo        entity.PageInfo         `json:"page_info"`
	OrderPharmacies []OrderPharmacyResponse `json:"order_pharmacies"`
}

type AllOrderPharmaciesSummaryResponse struct {
	AllCount       int64 `json:"all_count"`
	UnpaidCount    int64 `json:"unpaid_count"`
	ApprovalCount  int64 `json:"approval_count"`
	PendingCount   int64 `json:"pending_count"`
	SentCount      int64 `json:"sent_count"`
	ConfirmedCount int64 `json:"confirmed_count"`
	CanceledCount  int64 `json:"canceled_count"`
}

type OrderPharmacyResponse struct {
	OrderPharmacyId       int64               `json:"order_pharmacy_id"`
	Address               string              `json:"address,omitempty"`
	OrderStatusId         int64               `json:"order_status_id"`
	SubtotalAmount        decimal.Decimal     `json:"subtotal_amount"`
	DeliveryFee           decimal.Decimal     `json:"delivery_fee"`
	PharmacyName          string              `json:"pharmacy_name"`
	PharmacistPhoneNumber string              `json:"pharmacist_phone_number,omitempty"`
	PharmacyManagerEmail  string              `json:"pharmacy_manager_email,omitempty"`
	CourierName           string              `json:"courier_name,omitempty"`
	ProfilePicture        string              `json:"profile_picture"`
	OrderItems            []OrderItemResponse `json:"order_items,omitempty"`
	OrderItemsCount       int64               `json:"order_items_count,omitempty"`
	FirstOrderItem        *OrderItemResponse  `json:"first_order_item,omitempty"`
	UpdatedAt             time.Time           `json:"updated_at,omitempty"`
	CreatedAt             time.Time           `json:"created_at,omitempty"`
}

type OrderItemResponse struct {
	Id        int64           `json:"id,omitempty"`
	DrugName  string          `json:"drug_name"`
	DrugPrice decimal.Decimal `json:"drug_price"`
	DrugUnit  string          `json:"drug_unit"`
	Quantity  int             `json:"quantity"`
	DrugImage string          `json:"drug_image"`
}

func ConvertToOrderResponse(order entity.Order) OrderResponse {
	return OrderResponse{
		Id:             order.Id,
		Address:        order.Address,
		PaymentProof:   order.PaymentProof,
		TotalAmount:    order.TotalAmount,
		PharmacyOrders: ConvertToAllOrderPharmaciesResponse(order.OrderPharmacies),
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}
}

func ConvertToOrderPharmacyResponse(orderPharmacy entity.OrderPharmacy, orderItems []entity.OrderItem) *OrderPharmacyResponse {
	orderPharmacyResponse := &OrderPharmacyResponse{
		OrderPharmacyId:       orderPharmacy.Id,
		Address:               orderPharmacy.Address,
		OrderStatusId:         orderPharmacy.OrderStatusId,
		SubtotalAmount:        orderPharmacy.SubtotalAmount,
		DeliveryFee:           orderPharmacy.DeliveryFee,
		PharmacyName:          orderPharmacy.PharmacyName,
		PharmacistPhoneNumber: orderPharmacy.PharmacistPhoneNumber,
		PharmacyManagerEmail:  orderPharmacy.PharmacyManagerEmail,
		CourierName:           orderPharmacy.CourierName,
		ProfilePicture:        orderPharmacy.ProfilePicture,
		OrderItems:            ConvertToAllOrderItemsResponse(orderItems),
		OrderItemsCount:       orderPharmacy.OrderItemsCount,
		FirstOrderItem:        nil,
		UpdatedAt:             orderPharmacy.UpdatedAt,
		CreatedAt:             orderPharmacy.CreatedAt,
	}

	if orderPharmacy.FirstOrderItem.DrugName != "" {
		orderPharmacyResponse.FirstOrderItem = ConvertToOrderItemResponse(orderPharmacy.FirstOrderItem)
	}

	return orderPharmacyResponse
}

func ConvertToOrderItemResponse(orderItem entity.OrderItem) *OrderItemResponse {
	return &OrderItemResponse{
		Id:        orderItem.Id,
		DrugName:  orderItem.DrugName,
		DrugPrice: orderItem.DrugPrice,
		DrugUnit:  orderItem.DrugUnit,
		Quantity:  orderItem.Quantity,
		DrugImage: orderItem.DrugImage,
	}
}

func ConvertToAllOrderPharmaciesResponse(orderPharmacies []entity.OrderPharmacy) []OrderPharmacyResponse {
	orderPharmaciesResponse := []OrderPharmacyResponse{}

	for _, orderPharmacy := range orderPharmacies {
		orderPharmaciesResponse = append(orderPharmaciesResponse, *ConvertToOrderPharmacyResponse(orderPharmacy, nil))
	}

	return orderPharmaciesResponse
}

func ConvertToAllOrderPharmaciesResponseWithPageInfoAndPointer(orderPharmacies []*entity.OrderPharmacy, pageInfo entity.PageInfo) *AllOrderPharmaciesResponse {
	orderPharmaciesResponse := []OrderPharmacyResponse{}

	for _, orderPharmacy := range orderPharmacies {
		orderPharmaciesResponse = append(orderPharmaciesResponse, *ConvertToOrderPharmacyResponse(*orderPharmacy, orderPharmacy.OrderItems))
	}

	return &AllOrderPharmaciesResponse{OrderPharmacies: orderPharmaciesResponse, PageInfo: pageInfo}
}

func ConvertToAllOrderPharmaciesResponseWithPageInfo(orderPharmacies []entity.OrderPharmacy, pageInfo entity.PageInfo) *AllOrderPharmaciesResponse {
	orderPharmaciesResponse := []OrderPharmacyResponse{}

	for _, orderPharmacy := range orderPharmacies {
		orderPharmaciesResponse = append(orderPharmaciesResponse, *ConvertToOrderPharmacyResponse(orderPharmacy, nil))
	}

	return &AllOrderPharmaciesResponse{OrderPharmacies: orderPharmaciesResponse, PageInfo: pageInfo}
}

func ConvertToAllOrderItemsResponse(orderItems []entity.OrderItem) []OrderItemResponse {
	orderItemsResponse := []OrderItemResponse{}

	for _, orderItem := range orderItems {
		orderItemsResponse = append(orderItemsResponse, *ConvertToOrderItemResponse(orderItem))
	}

	return orderItemsResponse
}

func ConvertToAllOrdersResponse(orders []*entity.Order, pageInfo entity.PageInfo) *AllOrdersResponse {
	ordersResponse := []OrderResponse{}

	for _, order := range orders {
		ordersResponse = append(ordersResponse, ConvertToOrderResponse(*order))
	}

	return &AllOrdersResponse{Orders: ordersResponse, PageInfo: pageInfo}
}

type PharmacyCheckoutRequest struct {
	PharmacyId        int64   `json:"pharmacy_id" validate:"required"`
	PharmacyCourierId int64   `json:"pharmacy_courier_id" validate:"required"`
	DeliveryFee       int     `json:"delivery_fee" validate:"required"`
	Subtotal          int     `json:"subtotal_amount" validate:"required"`
	CartItemIds       []int64 `json:"cart_items" validate:"required"`
}

type PharmacyCheckoutFromPrescriptionRequest struct {
	PharmacyId        int64                  `json:"pharmacy_id" validate:"required"`
	PharmacyCourierId int64                  `json:"pharmacy_courier_id" validate:"required"`
	DeliveryFee       int                    `json:"delivery_fee" validate:"required"`
	Subtotal          int                    `json:"subtotal_amount" validate:"required"`
	PharmacyDrugs     []PharmacyDrugQuantity `json:"pharmacy_drugs" validate:"required,min=1"`
}

type OrderCheckoutRequest struct {
	AccountId   int64
	Address     string                    `json:"address" binding:"required"`
	TotalAmount int                       `json:"total_amount" binding:"required"`
	Pharmacies  []PharmacyCheckoutRequest `json:"pharmacies" binding:"required"`
}

type PharmacyDrugQuantity struct {
	PharmacyDrugId int64 `json:"pharmacy_drug_id" validate:"required,gte=1"`
	Quantity       int   `json:"quantity" validate:"required,gte=1"`
}

type OrderCheckoutResponse struct {
	OrderId int64 `json:"order_id"`
}

type OrderChangeStatusRequest struct {
	StatusId int64 `json:"status_id" binding:"required"`
}

func ConvertPrescriptionCheckoutRequest(request CheckoutFromPrescriptionRequest) OrderCheckoutRequest {
	return OrderCheckoutRequest{
		AccountId:   request.AccountId,
		TotalAmount: request.TotalAmount,
		Address:     request.Address,
		Pharmacies:  ConvertPharmacyCheckoutFromPrescriptionRequestList(request.Pharmacies),
	}
}

func ConvertPharmacyCheckoutFromPrescriptionRequestList(requestList []PharmacyCheckoutFromPrescriptionRequest) []PharmacyCheckoutRequest {
	var resultList []PharmacyCheckoutRequest

	for _, request := range requestList {
		resultList = append(resultList, ConvertPharmacyCheckoutFromPrescriptionRequest(request))
	}

	return resultList
}

func ConvertPharmacyCheckoutFromPrescriptionRequest(request PharmacyCheckoutFromPrescriptionRequest) PharmacyCheckoutRequest {
	return PharmacyCheckoutRequest{
		PharmacyId:        request.PharmacyId,
		PharmacyCourierId: request.PharmacyCourierId,
		DeliveryFee:       request.DeliveryFee,
		Subtotal:          request.Subtotal,
	}
}
