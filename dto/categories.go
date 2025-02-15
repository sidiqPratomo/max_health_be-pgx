package dto

import "github.com/sidiqPratomo/max-health-backend/entity"

type CategoryRequest struct {
	Name string `json:"category_name" validate:"required"`
}

type CategoryResponse struct {
	Id   int64  `json:"category_id"`
	Name string `json:"category_name"`
	Url  string `json:"category_url"`
}

func ConvertCategoryToCategoryResponse(category entity.DrugCategory) CategoryResponse {
	return CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
		Url:  category.Url,
	}
}

func ConvertCategoriesToCategoriesResponse(categories []entity.DrugCategory) []CategoryResponse {
	data := []CategoryResponse{}
	for _, category := range categories {
		data = append(data, ConvertCategoryToCategoryResponse(category))
	}
	return data
}

func ConvertCategoryRequestToEntity(categoryRequest CategoryRequest) entity.DrugCategory {
	return entity.DrugCategory{
		Name: categoryRequest.Name,
	}
}
