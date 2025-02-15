package dto

type StockChangeQuery struct {
	PharmacyId *int64 `form:"pharmacy-id"`
}

type StockChangeResponse struct {
	PharmacyName    string `json:"pharmacy_name"`
	PharmacyAddress string `json:"pharmacy_address"`
	DrugImage       string `json:"drug_url"`
	DrugName        string `json:"drug_name"`
	FinalStock      int    `json:"final_stock"`
	Change          int    `json:"stock_change"`
	Description     string `json:"description"`
}
