package usecase

import (
	"context"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type CategoryUsecase interface {
	GetAllCategories(ctx context.Context) ([]entity.DrugCategory, error)
	DeleteOneCategoryById(ctx context.Context, categoryId int64) error
	AddOneCategory(ctx context.Context, newCategory entity.DrugCategory, file multipart.File, fileHeader multipart.FileHeader) error
	UpdateOneCategory(ctx context.Context, updatedCategory entity.DrugCategory, file multipart.File, fileHeader *multipart.FileHeader) error
}

type categoryUsecaseImpl struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryUsecaseImpl(categoryRepository repository.CategoryRepository) categoryUsecaseImpl {
	return categoryUsecaseImpl{
		categoryRepository: categoryRepository,
	}
}

func (u *categoryUsecaseImpl) GetAllCategories(ctx context.Context) ([]entity.DrugCategory, error) {
	categories, err := u.categoryRepository.FindAllCategories((ctx))
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (u *categoryUsecaseImpl) DeleteOneCategoryById(ctx context.Context, categoryId int64) error {
	category, err := u.categoryRepository.FindOneCategoryById(ctx, categoryId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if category == nil {
		return apperror.CategoryNotFoundError()
	}

	if strings.Split(category.Url, "/")[1] == "res.cloudinary.com" {
		util.DeleteInCloudinary(category.Url)
	}

	err = u.categoryRepository.DeleteOneCategoryById(ctx, categoryId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *categoryUsecaseImpl) AddOneCategory(ctx context.Context, newCategory entity.DrugCategory, file multipart.File, fileHeader multipart.FileHeader) error {
	category, err := u.categoryRepository.FindOneCategoryByName(ctx, newCategory.Name)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if category != nil {
		return apperror.CategoryNotUniqueError()
	}

	filePath, _, err := util.ValidateFile(fileHeader, appconstant.CategoryPicturesUrl, []string{"png", "jpg", "jpeg"}, 2000000)
	if err != nil {
		return apperror.NewAppError(http.StatusBadRequest, err, err.Error())
	}

	imageUrl, err := util.UploadToCloudinary(file, *filePath)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	newCategory.Url = imageUrl

	err = u.categoryRepository.PostOneCategory(ctx, newCategory)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *categoryUsecaseImpl) UpdateOneCategory(ctx context.Context, updatedCategory entity.DrugCategory, file multipart.File, fileHeader *multipart.FileHeader) error {
	category, err := u.categoryRepository.FindSimilarCategory(ctx, updatedCategory)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if category != nil {
		return apperror.CategoryNotUniqueError()
	}

	dbCategory, err := u.categoryRepository.FindOneCategoryById(ctx, updatedCategory.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if dbCategory == nil {
		return apperror.CategoryNotFoundError()
	}

	if file != nil {
		filePath, _, err := util.ValidateFile(*fileHeader, appconstant.CategoryPicturesUrl, []string{"png", "jpg", "jpeg"}, 2000000)
		if err != nil {
			return apperror.NewAppError(http.StatusBadRequest, err, err.Error())
		}

		imageUrl, err := util.UploadToCloudinary(file, *filePath)
		if err != nil {
			return apperror.InternalServerError(err)
		}
		updatedCategory.Url = imageUrl
		if strings.Split(dbCategory.Url, "/")[1] == "res.cloudinary.com" {
			util.DeleteInCloudinary(dbCategory.Url)
		}
	} else {
		updatedCategory.Url = dbCategory.Url
	}

	err = u.categoryRepository.UpdateOneCategoryById(ctx, updatedCategory)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
