package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/shopspring/decimal"
)

type DrugUsecase interface {
	GetPharmacyDrugByDrugId(ctx context.Context, drugId int64, latitude, longitude, page, limit string) (*dto.DrugDetailResponse, error)
	GetAllDrugsForListing(ctx context.Context, query *util.ValidatedGetProductQuery) ([]entity.DrugListing, *entity.PageInfo, error)
	UpdateOneDrug(ctx context.Context, drugId int64, drugRequest dto.UpdateDrugRequest, file multipart.File, fileHeader *multipart.FileHeader) error
	GetAllDrugs(ctx context.Context, query *util.ValidatedGetDrugAdminQuery) (*dto.AllDrugsResponse, error)
	GetOneDrugByDrugId(ctx context.Context, drugId int64) (*dto.DrugResponse, error)
	CreateOneDrug(ctx context.Context, drugRequest dto.CreateDrugRequest, file multipart.File, fileHeader *multipart.FileHeader) error
	DeleteOneDrug(ctx context.Context, drugId int64) error
	GetDrugsByPharmacyId(ctx context.Context, pharmacyId string, limit string, page string, search string) (*dto.PharmacyDrugsByPharmacyResponse, error)
	UpdateDrugsByPharmacyDrugId(ctx context.Context, pharmacyDrugId int64, stock int, price decimal.Decimal) error
	DeleteDrugsByPharmacyDrugId(ctx context.Context, pharmacyDrugId int64) error
	AddDrugsByPharmacyDrugId(ctx context.Context, pharmacyId int64, drugId int64, stock int, price decimal.Decimal) error
	GetPossibleStockMutation(ctx context.Context, pharmacyDrugId int64) ([]dto.PharmacyDrugMutationsResponse, error)
	PostStockMutation(ctx context.Context, req dto.PostStockMutationRequest) error
}

type drugUsecaseImpl struct {
	transaction                  repository.Transaction
	drugRepository               repository.DrugRepository
	pharmacyDrugRepository       repository.PharmacyDrugRepository
	categoryRepository           repository.CategoryRepository
	drugClassificationRepository repository.DrugClassificationRepository
	drugFormRepository           repository.DrugFormRepository
	pharmacyRepository           repository.PharmacyRepository
}

func NewDrugUsecaseImpl(transaction repository.Transaction, drugRepository repository.DrugRepository, pharmacyDrugRepository repository.PharmacyDrugRepository, drugClassificationRepository repository.DrugClassificationRepository, drugFormRepository repository.DrugFormRepository, categoryRepository repository.CategoryRepository, pharmacyRepository repository.PharmacyRepository) drugUsecaseImpl {
	return drugUsecaseImpl{
		transaction:                  transaction,
		drugRepository:               drugRepository,
		pharmacyDrugRepository:       pharmacyDrugRepository,
		drugClassificationRepository: drugClassificationRepository,
		drugFormRepository:           drugFormRepository,
		categoryRepository:           categoryRepository,
		pharmacyRepository:           pharmacyRepository,
	}
}

func (u *drugUsecaseImpl) GetPharmacyDrugByDrugId(ctx context.Context, drugId int64, latitude, longitude, page, limit string) (*dto.DrugDetailResponse, error) {
	latitudeFloat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return nil, apperror.CoordinateInvalidError()
	}

	longitudeFloat, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return nil, apperror.CoordinateInvalidError()
	}

	limitInt, offset, err := util.CheckPharmacyDrugPagination(page, limit)
	if err != nil {
		return nil, err
	}

	drugDetail, err := u.drugRepository.GetOneActiveDrugById(ctx, drugId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if drugDetail == nil {
		return nil, apperror.DrugNotFoundError()
	}

	pharmacyDrugList, err := u.pharmacyDrugRepository.GetPharmacyDrugsByDrugId(ctx, drugId, latitudeFloat, longitudeFloat, limitInt, offset)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	pharmacyDrugListResponse := dto.ConvertToDrugDetailResponse(*drugDetail, pharmacyDrugList)

	return &pharmacyDrugListResponse, err
}

func (u *drugUsecaseImpl) GetAllDrugsForListing(ctx context.Context, query *util.ValidatedGetProductQuery) ([]entity.DrugListing, *entity.PageInfo, error) {
	if query.Category != nil {
		category, err := u.categoryRepository.FindOneCategoryById(ctx, *query.Category)
		if err != nil {
			return nil, nil, apperror.InternalServerError(err)
		}
		if category == nil {
			return nil, nil, apperror.CategoryNotFoundError()
		}
	}

	drugList, pageInfo, err := u.pharmacyDrugRepository.GetProductListing(ctx, query)
	if err != nil {
		return nil, nil, apperror.InternalServerError(err)
	}
	return drugList, pageInfo, nil
}

func (u *drugUsecaseImpl) UpdateOneDrug(ctx context.Context, drugId int64, drugRequest dto.UpdateDrugRequest, file multipart.File, fileHeader *multipart.FileHeader) error {
	drug := dto.ConvertToDrug(drugRequest)
	drug.Id = drugId

	existingDrug, err := u.drugRepository.GetOneDrugById(ctx, drug.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if existingDrug == nil {
		return apperror.DrugNotFoundError()
	}

	sameNameDrug, err := u.drugRepository.GetDrugByName(ctx, drug.Name)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if sameNameDrug != nil && sameNameDrug.Id != drug.Id {
		if drug.Manufacture == sameNameDrug.Manufacture && drug.Content == sameNameDrug.Content && drug.GenericName == sameNameDrug.GenericName {
			return apperror.DrugNameAlreadyExistError()
		}
	}

	if file != nil {
		filePath, _, err := util.ValidateFile(*fileHeader, appconstant.DrugPicturesUrl, []string{"png"}, 500000)
		if err != nil {
			return apperror.NewAppError(http.StatusBadRequest, err, err.Error())
		}

		imageUrl, err := util.UploadToCloudinary(file, *filePath)
		if err != nil {
			return apperror.InternalServerError(err)
		}

		drug.Image = imageUrl
		if strings.Contains(existingDrug.Image, "res.cloudinary.com") {
			util.DeleteInCloudinary(existingDrug.Image)
		}
	} else {
		drug.Image = existingDrug.Image
	}

	err = u.drugRepository.UpdateOneDrug(ctx, drug)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *drugUsecaseImpl) GetAllDrugs(ctx context.Context, query *util.ValidatedGetDrugAdminQuery) (*dto.AllDrugsResponse, error) {
	drugs, pageInfo, err := u.drugRepository.GetAllDrugs(ctx, *query)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	res := &dto.AllDrugsResponse{Drugs: drugs, PageInfo: *pageInfo}

	return res, nil
}

func (u *drugUsecaseImpl) GetOneDrugByDrugId(ctx context.Context, drugId int64) (*dto.DrugResponse, error) {
	drug, err := u.drugRepository.GetOneDrugById(ctx, drugId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if drug == nil {
		return nil, apperror.DrugNotFoundError()
	}

	drugResponse := dto.ConvertToDrugResponse(*drug)

	return &drugResponse, err
}

func (u *drugUsecaseImpl) CreateOneDrug(ctx context.Context, createDrugRequest dto.CreateDrugRequest, file multipart.File, fileHeader *multipart.FileHeader) error {
	drug := dto.CreateDrugRequestToDrug(createDrugRequest)

	existingDrugId, err := u.drugRepository.GetDrugIdByName(ctx, drug.Name)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if existingDrugId != nil {
		return apperror.DrugNameAlreadyExistError()
	}

	classification, err := u.drugClassificationRepository.FindOneById(ctx, drug.Classification.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if classification == nil {
		return apperror.ClassificationNotFoundError()
	}

	drugForm, err := u.drugFormRepository.FindOneById(ctx, drug.Form.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if drugForm == nil {
		return apperror.DrugFormNotFoundError()
	}

	category, err := u.categoryRepository.FindOneCategoryById(ctx, drug.Category.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if category == nil {
		return apperror.CategoryNotFoundError()
	}

	filePath, _, err := util.ValidateFile(*fileHeader, appconstant.DrugPicturesUrl, []string{"png", "jpg", "jpeg"}, 2000000)
	if err != nil {
		return apperror.NewAppError(http.StatusBadRequest, err, err.Error())
	}

	imageUrl, err := util.UploadToCloudinary(file, *filePath)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	drug.Image = imageUrl

	if err = u.drugRepository.CreateOneDrug(ctx, drug); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *drugUsecaseImpl) DeleteOneDrug(ctx context.Context, drugId int64) error {
	drug, err := u.drugRepository.GetOneDrugById(ctx, drugId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if drug == nil {
		return apperror.DrugNotFoundError()
	}

	if err = u.drugRepository.DeleteOneDrug(ctx, drugId); err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *drugUsecaseImpl) GetDrugsByPharmacyId(ctx context.Context, pharmacyId string, limit string, page string, search string) (*dto.PharmacyDrugsByPharmacyResponse, error) {
	pharmacyIdInt, err := strconv.Atoi(pharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 12
	}

	offset := (pageInt - 1) * limitInt
	if offset < 0 {
		offset = 0
	}

	drugs, pageInfo, err := u.drugRepository.GetDrugsByPharmacyId(ctx, int64(pharmacyIdInt), limit, offset, search)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if drugs == nil {
		return nil, apperror.DrugNotFoundError()
	}

	getPharmacyDrugByPharmacyId := []dto.PharmacyDrugByPharmacyDTO{}

	for _, drug := range drugs {
		pharmacyDto := dto.PharmacyDrugByPharmacyDTO{
			Id:    drug.Id,
			Price: drug.Price,
			Stock: drug.Stock,
			Drug: dto.DrugResponse{
				Id:          drug.Drug.Id,
				Name:        drug.Drug.Name,
				GenericName: drug.Drug.GenericName,
				Content:     drug.Drug.Content,
				Manufacture: drug.Drug.Manufacture,
				Description: drug.Drug.Description,
				Classification: dto.DrugClassification{
					Id:   drug.Drug.Classification.Id,
					Name: drug.Drug.Classification.Name,
				},
				Form: dto.DrugForm{
					Id:   drug.Drug.Form.Id,
					Name: drug.Drug.Form.Name,
				},
				UnitInPack:  drug.Drug.UnitInPack,
				SellingUnit: drug.Drug.SellingUnit,
				Weight:      drug.Drug.Weight,
				Height:      drug.Drug.Height,
				Length:      drug.Drug.Length,
				Width:       drug.Drug.Width,
				Image:       drug.Drug.Image,
				Category: dto.DrugCategory{
					Id:   drug.Drug.Category.Id,
					Name: drug.Drug.Category.Name,
					Url:  drug.Drug.Category.Url,
				},
				IsActive:               drug.Drug.IsActive,
				IsPreScriptionRequired: drug.Drug.IsPrescriptionRequired,
			},
		}

		getPharmacyDrugByPharmacyId = append(getPharmacyDrugByPharmacyId, pharmacyDto)
	}

	getDrugsByPharmacyResponse := dto.PharmacyDrugsByPharmacyResponse{
		Drugs:    getPharmacyDrugByPharmacyId,
		PageInfo: *pageInfo,
	}

	return &getDrugsByPharmacyResponse, nil
}

func (u *drugUsecaseImpl) UpdateDrugsByPharmacyDrugId(ctx context.Context, pharmacyDrugId int64, stock int, price decimal.Decimal) error {
	if stock < 0 {
		return apperror.BadRequestError(errors.New("stock cannot be less than 0"))
	}
	if price.Cmp(decimal.NewFromInt(500)) < 0 {
		return apperror.BadRequestError(errors.New("price cannot be less than 500"))
	}
	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	pharmacyDrugRepo := tx.PharmacyDrugRepo()
	stockChangeRepo := tx.StockChangeRepo()
	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()
	pharmacyDrug, err := pharmacyDrugRepo.GetPharmacyDrugByIdForUpdate(ctx, pharmacyDrugId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if pharmacyDrug == nil {
		err = apperror.DrugNotFoundError()
		return err
	}
	err = pharmacyDrugRepo.UpdatePharmacyDrugStockPrice(ctx, pharmacyDrugId, stock, price)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	stockChange := entity.StockChange{PharmacyDrugId: pharmacyDrug.Id, FinalStock: stock, Amount: stock - pharmacyDrug.Stock,
		Description: "updated by manager",}
	err = stockChangeRepo.PostStockChangesFromUpdate(ctx, []entity.StockChange{stockChange})
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *drugUsecaseImpl) DeleteDrugsByPharmacyDrugId(ctx context.Context, pharmacyDrugId int64) error {
	err := u.pharmacyDrugRepository.DeletePharmacyDrug(ctx, pharmacyDrugId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *drugUsecaseImpl) AddDrugsByPharmacyDrugId(ctx context.Context, pharmacyId int64, drugId int64, stock int, price decimal.Decimal) error {
	if stock < 0 {
		return apperror.BadRequestError(errors.New("stock cannot be less than 0"))
	}
	if price.Cmp(decimal.NewFromInt(500)) < 0 {
		return apperror.BadRequestError(errors.New("price cannot be less than 500"))
	}

	pharmacy, err := u.pharmacyRepository.GetOnePharmacyByPharmacyId(ctx, pharmacyId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if pharmacy == nil {
		return apperror.BadRequestError(errors.New("pharmacy id doesn't exist"))
	}

	drug, err := u.drugRepository.GetDrugById(ctx, drugId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if drug == nil {
		return apperror.BadRequestError(errors.New("drug id doesn't exist"))
	}

	err = u.pharmacyDrugRepository.AddPharmacyDrug(ctx, pharmacyId, drugId, stock, price)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *drugUsecaseImpl) GetPossibleStockMutation(ctx context.Context, pharmacyDrugId int64) ([]dto.PharmacyDrugMutationsResponse, error) {
	currentDrug, err := u.pharmacyDrugRepository.GetPharmacyDrugById(ctx, pharmacyDrugId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if currentDrug == nil {
		return nil, apperror.DrugNotFoundError()
	}

	pharmacyDrugs, err := u.pharmacyDrugRepository.GetPossibleStockMutation(ctx, pharmacyDrugId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToMutationPharmacyDrugs(pharmacyDrugs), nil
}

func (u *drugUsecaseImpl) PostStockMutation(ctx context.Context, req dto.PostStockMutationRequest) error {
	if req.RecipientPharmacyDrugId == req.SenderPharmacyDrugId {
		return apperror.DuplicatePharmacyDrugIdError()
	}
	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	pharmacyDrugRepo := tx.PharmacyDrugRepo()
	stockMutationRepo := tx.StockMutationRepo()
	stockChangeRepo := tx.StockChangeRepo()
	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	recipientDrug, err := pharmacyDrugRepo.GetPharmacyDrugByIdForUpdate(ctx, req.RecipientPharmacyDrugId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	senderDrug, err := pharmacyDrugRepo.GetPharmacyDrugByIdForUpdate(ctx, req.SenderPharmacyDrugId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if recipientDrug == nil || senderDrug == nil {
		err = apperror.DrugNotFoundError()
		return err
	}
	if recipientDrug.DrugId != senderDrug.DrugId {
		err = apperror.InvalidStockMutationRequestError()
		return err
	}

	if senderDrug.Stock < req.Quantity {
		err = apperror.InsufficientStockError()
		return err
	}

	err = pharmacyDrugRepo.UpdatePharmacyDrugStockPrice(ctx, req.RecipientPharmacyDrugId, recipientDrug.Stock+req.Quantity, recipientDrug.Price)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	err = pharmacyDrugRepo.UpdatePharmacyDrugStockPrice(ctx, req.SenderPharmacyDrugId, senderDrug.Stock-req.Quantity, senderDrug.Price)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	stockMutation := entity.PossibleStockMutation{DrugId: recipientDrug.DrugId, OriginalPharmacy: recipientDrug.PharmacyId,
		AlternativePharmacy: senderDrug.PharmacyId, AlternativeStock: req.Quantity}
	err = stockMutationRepo.PostStockMutations(ctx, []entity.PossibleStockMutation{stockMutation})
	if err != nil {
		return apperror.InternalServerError(err)
	}

	recipientStockChange := entity.StockChange{PharmacyDrugId: req.RecipientPharmacyDrugId, FinalStock: recipientDrug.Stock + req.Quantity,
		Amount: req.Quantity}
	senderStockChange := entity.StockChange{PharmacyDrugId: req.SenderPharmacyDrugId, FinalStock: senderDrug.Stock - req.Quantity,
		Amount: req.Quantity * -1}
	err = stockChangeRepo.PostStockChangesFromMutation(ctx, []entity.StockChange{recipientStockChange, senderStockChange})
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
