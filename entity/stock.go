package entity

import "time"

type StockChange struct {
	Id             int64
	PharmacyDrugId int64
	FinalStock     int
	Amount         int
	Description    string
	CreatedAt      time.Time
}

type StockMutationRequest struct {
	Id                  int64
	PharmacyRequesterId int64
	PharmacyTargetId    int64
	DrugId              int64
	Stock               int
	StatusId            int64
	CreatedAt           time.Time
}

type StockRequestStatus struct {
	Id   int64
	Name string
}

type PossibleStockMutation struct {
	CartItemId int64
	CartQuantity int
	DrugId int64
	OriginalPharmacyDrug int64
	OriginalPharmacy int64
	OriginalStock int
	AlternativePharmacyDrug int64
	AlternativePharmacy int64
	AlternativeStock int
}
