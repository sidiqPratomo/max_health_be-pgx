package dto

import "github.com/sidiqPratomo/max-health-backend/entity"

type DrugCategorySalesVolumeRevenueResponse struct {
	DrugCategoryId   int64  `json:"drug_category_id"`
	DrugCategoryName string `json:"drug_category_name"`
	SalesVolume      int64  `json:"sales_volume"`
	Revenue          int64  `json:"revenue"`
}

type AllDrugCategorySalesVolumeRevenueResponse struct {
	Report []DrugCategorySalesVolumeRevenueResponse `json:"report"`
}

func ConvertToDrugCategorySalesVolumeRevenueResponse(report entity.DrugCategorySalesVolumeRevenue) DrugCategorySalesVolumeRevenueResponse {
	return DrugCategorySalesVolumeRevenueResponse{
		DrugCategoryId:   report.DrugCategoryId,
		DrugCategoryName: report.DrugCategoryName,
		SalesVolume:      report.SalesVolume,
		Revenue:          report.Revenue,
	}
}

func ConvertToAllDrugCategorySalesVolumeRevenueResponse(reports []entity.DrugCategorySalesVolumeRevenue) *AllDrugCategorySalesVolumeRevenueResponse {
	reportResponse := []DrugCategorySalesVolumeRevenueResponse{}

	for _, report := range reports {
		reportResponse = append(reportResponse, ConvertToDrugCategorySalesVolumeRevenueResponse(report))
	}

	return &AllDrugCategorySalesVolumeRevenueResponse{
		Report: reportResponse,
	}
}

type DrugSalesVolumeRevenueResponse struct {
	DrugId      int64  `json:"drug_id"`
	DrugName    string `json:"drug_name"`
	SalesVolume int64  `json:"sales_volume"`
	Revenue     int64  `json:"revenue"`
}

type AllDrugSalesVolumeRevenueResponse struct {
	Report []DrugSalesVolumeRevenueResponse `json:"report"`
}

func ConvertToDrugSalesVolumeRevenueResponse(report entity.DrugSalesVolumeRevenue) DrugSalesVolumeRevenueResponse {
	return DrugSalesVolumeRevenueResponse{
		DrugId:      report.DrugId,
		DrugName:    report.DrugName,
		SalesVolume: report.SalesVolume,
		Revenue:     report.Revenue,
	}
}

func ConvertToAllDrugSalesVolumeRevenueResponse(reports []entity.DrugSalesVolumeRevenue) *AllDrugSalesVolumeRevenueResponse {
	reportResponse := []DrugSalesVolumeRevenueResponse{}

	for _, report := range reports {
		reportResponse = append(reportResponse, ConvertToDrugSalesVolumeRevenueResponse(report))
	}

	return &AllDrugSalesVolumeRevenueResponse{
		Report: reportResponse,
	}
}
