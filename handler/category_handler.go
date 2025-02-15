package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase) CategoryHandler {
	return CategoryHandler{
		categoryUsecase: categoryUsecase,
	}
}

func (h *CategoryHandler) GetAllCategories(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	categories, err := h.categoryUsecase.GetAllCategories(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	resCategories := dto.ConvertCategoriesToCategoriesResponse(categories)
	util.ResponseOK(ctx, resCategories)
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	paramCategoryId := ctx.Param(appconstant.CategoryIdString)
	categoryId, _ := strconv.Atoi(paramCategoryId)
	err := h.categoryUsecase.DeleteOneCategoryById(ctx, int64(categoryId))
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseOK(ctx, nil)
}

func (h *CategoryHandler) AddOneCategory(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var categoryRequest dto.CategoryRequest
	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.NewAppError(http.StatusBadRequest, errors.New("data is empty"), "data is empty"))
		return
	}

	err := json.Unmarshal([]byte(data), &categoryRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(categoryRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		if file == nil {
			ctx.Error(apperror.NewAppError(http.StatusBadRequest, errors.New("file not attached"), "file not attached"))
			return
		}
		ctx.Error(err)
		return
	}

	category := dto.ConvertCategoryRequestToEntity(categoryRequest)
	err = h.categoryUsecase.AddOneCategory(ctx.Request.Context(), category, file, *fileHeader)
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseCreated(ctx, nil)
}

func (h *CategoryHandler) UpdateOneCategory(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var categoryRequest dto.CategoryRequest
	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.NewAppError(http.StatusBadRequest, errors.New("data is empty"), "data is empty"))
		return
	}

	err := json.Unmarshal([]byte(data), &categoryRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(categoryRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		if file != nil {
			ctx.Error(err)
			return
		}
	}

	category := dto.ConvertCategoryRequestToEntity(categoryRequest)
	paramCategoryId := ctx.Param(appconstant.CategoryIdString)
	categoryId, _ := strconv.Atoi(paramCategoryId)
	category.Id = int64(categoryId)

	err = h.categoryUsecase.UpdateOneCategory(ctx, category, file, fileHeader)
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseOK(ctx, nil)
}
