package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderItem struct {
	Id              int64
	OrderPharmacyId int64
	DrugName        string
	DrugPrice       decimal.Decimal
	DrugUnit        string
	Quantity        int
	DrugImage       string
}

type OrderPharmacy struct {
	Id                    int64
	Address               string
	OrderId               int64
	UserId                int64
	OrderStatusId         int64
	PharmacyCourierId     int64
	SubtotalAmount        decimal.Decimal
	DeliveryFee           decimal.Decimal
	PharmacyName          string
	PharmacistPhoneNumber string
	CourierName           string
	ProfilePicture        string
	PharmacyManagerEmail  string
	OrderItemsCount       int64
	OrderItems            []OrderItem
	FirstOrderItem        OrderItem
	UpdatedAt             time.Time
	CreatedAt             time.Time
}

type OrderPharmacyForCheckout struct {
	Id        int64
	CartItems []CartItemForCheckout
}

type OrderPharmacySummary struct {
	AllCount       int64
	UnpaidCount    int64
	ApprovalCount  int64
	PendingCount   int64
	SentCount      int64
	ConfirmedCount int64
	CanceledCount  int64
}

type Order struct {
	Id              int64
	UserId          int64
	Address         string
	PaymentProof    string
	TotalAmount     decimal.Decimal
	OrderPharmacies []OrderPharmacy
	ExpiredAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderStatus struct {
	Id   int64
	Name string
}

type DrugCategorySalesVolumeRevenue struct {
	DrugCategoryId   int64
	DrugCategoryName string
	SalesVolume      int64
	Revenue          int64
}

type DrugSalesVolumeRevenue struct {
	DrugId      int64
	DrugName    string
	SalesVolume int64
	Revenue     int64
}
