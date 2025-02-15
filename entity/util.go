package entity

type PageInfo struct {
	PageCount int `json:"page_count"`
	ItemCount int `json:"item_count"`
	Page      int `json:"page"`
}
