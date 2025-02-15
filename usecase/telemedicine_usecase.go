package usecase

import (
	"context"
	"math"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/shopspring/decimal"
)

type TelemedicineUsecase interface {
	UserCreateRoom(ctx context.Context, userAccountId, doctorAccountId int64) (*int64, error)
	DoctorJoinRoom(ctx context.Context, doctorAccountId, roomId int64) error
	PostOneMessage(ctx context.Context, accountId int64, postOneMessageRequest dto.PostOneMessageRequest, file multipart.File, fileHeader *multipart.FileHeader) (*dto.Chat, error)
	Listen(ctx context.Context, accountId int64, roomId int64) (*dto.Chat, error)
	findListenerBySender(sender entity.Participant) bool
	findListener(listener entity.Participant) bool
	addListener(listener entity.Participant)
	removeListener(listener entity.Participant)
	GetAllChat(ctx context.Context, accountId, roomId int64) (*dto.ChatRoom, error)
	GetAllChatRoomPreview(ctx context.Context, accountId int64, role string) ([]dto.ChatRoomPreview, error)
	DoctorGetChatRequest(ctx context.Context, accountId int64) ([]dto.ChatRoomPreview, error)
	SavePrescription(ctx context.Context, accountId, prescriptionId int64) error
	GetAllPrescriptions(ctx context.Context, accountId int64, limit, page string) (*dto.PrescriptionResponseList, error)
	PrepareForCheckout(ctx context.Context, accountId, prescriptionId int64, addressIdString string) (*dto.PreapareForCheckoutResponse, error)
	CheckoutFromPrescription(ctx context.Context, checkoutFromPrescriptionRequest dto.CheckoutFromPrescriptionRequest) (*int64, error)
	CloseChatRoom(ctx context.Context, userAccountId, roomId int64) error
}

type telemedicineUsecaseImpl struct {
	chatRoomRepository         repository.ChatRoomRepository
	chatRepository             repository.ChatRepository
	userRepository             repository.UserRepository
	doctorRepository           repository.DoctorRepository
	pharmacyDrugRepository     repository.PharmacyDrugRepository
	prescriptionDrugRepository repository.PrescriptionDrugRepository
	prescriptionRepository     repository.PrescriptionRepository
	cartRepository             repository.CartRepository
	orderRepository            repository.OrderRepository
	userAddressRepository      repository.UserAddressRepository
	pharmacyRepository         repository.PharmacyRepository
	chatChannel                map[int64]chan entity.Chat
	listeners                  []entity.Participant
	listenersLock              sync.Mutex
	abortChannel               chan entity.Participant
	transaction                repository.Transaction
}

func NewTelemedicineUsecaseImpl(chatRoomRepository repository.ChatRoomRepository, chatRepository repository.ChatRepository, userRepository repository.UserRepository, doctorRepository repository.DoctorRepository, pharmacyDrugRepository repository.PharmacyDrugRepository, prescriptionDrugRepository repository.PrescriptionDrugRepository, prescriptionRepository repository.PrescriptionRepository, cartRepository repository.CartRepository, orderRepository repository.OrderRepository, userAddressRepository repository.UserAddressRepository, pharmacyRepository repository.PharmacyRepository, transaction repository.Transaction) telemedicineUsecaseImpl {
	return telemedicineUsecaseImpl{
		chatRoomRepository:         chatRoomRepository,
		chatRepository:             chatRepository,
		userRepository:             userRepository,
		doctorRepository:           doctorRepository,
		pharmacyDrugRepository:     pharmacyDrugRepository,
		prescriptionDrugRepository: prescriptionDrugRepository,
		prescriptionRepository:     prescriptionRepository,
		cartRepository:             cartRepository,
		orderRepository:            orderRepository,
		userAddressRepository:      userAddressRepository,
		pharmacyRepository:         pharmacyRepository,
		chatChannel:                make(map[int64]chan entity.Chat),
		listeners:                  make([]entity.Participant, 0),
		listenersLock:              sync.Mutex{},
		abortChannel:               make(chan entity.Participant),
		transaction:                transaction,
	}
}

func (u *telemedicineUsecaseImpl) UserCreateRoom(ctx context.Context, userAccountId, doctorAccountId int64) (*int64, error) {
	user, err := u.userRepository.FindUserByAccountId(ctx, userAccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if user == nil {
		return nil, apperror.ForbiddenAction()
	}

	doctor, err := u.doctorRepository.FindDoctorByAccountId(ctx, doctorAccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if doctor == nil {
		return nil, apperror.ForbiddenAction()
	}

	room, err := u.chatRoomRepository.FindActiveChatRoom(ctx, userAccountId, doctorAccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if room != nil {
		return nil, apperror.OnGoingChatExistError()
	}

	roomId, err := u.chatRoomRepository.CreateOneRoom(ctx, userAccountId, doctorAccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	u.chatChannel[*roomId] = make(chan entity.Chat)

	return roomId, nil
}

func (u *telemedicineUsecaseImpl) DoctorJoinRoom(ctx context.Context, doctorAccountId, roomId int64) error {
	room, err := u.chatRoomRepository.FindChatRoomById(ctx, roomId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if room == nil {
		return apperror.BadRequestError(err)
	}
	if room.DoctorAccountId != doctorAccountId {
		return apperror.ForbiddenAction()
	}

	err = u.chatRoomRepository.StartChat(ctx, roomId, doctorAccountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *telemedicineUsecaseImpl) PostOneMessage(ctx context.Context, accountId int64, postOneMessageRequest dto.PostOneMessageRequest, file multipart.File, fileHeader *multipart.FileHeader) (*dto.Chat, error) {
	chat := dto.ConvertPostMessageRequestToChatEntity(postOneMessageRequest)
	chat.SenderAccountId = accountId

	chatRoom, err := u.chatRoomRepository.FindChatRoomById(ctx, chat.RoomId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if chatRoom == nil {
		return nil, apperror.ChatRoomNotFoundError()
	}
	if chatRoom.DoctorAccountId != accountId && chatRoom.UserAccountId != accountId {
		return nil, apperror.ForbiddenAction()
	}

	if chatRoom.ExpiredAt != nil {
		if chatRoom.ExpiredAt.Before(time.Now()) {
			return nil, apperror.RoomIsNowExpiredError()
		}
	}

	chat.RoomId = chatRoom.Id

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	prescriptionRepo := tx.PrescriptionRepository()
	prescriptionDrugRepo := tx.PrescriptionDrugRepository()
	chatRepo := tx.ChatRepository()

	if len(postOneMessageRequest.PrescriptionDrugs) > 0 {
		prescriptionId, err := prescriptionRepo.CreateOnePrescription(ctx, chatRoom.UserAccountId, chatRoom.DoctorAccountId)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		for _, prescriptionDrug := range chat.Prescription.PrescriptionDrugs {
			err := prescriptionDrugRepo.PostOnePrescriptionDrug(ctx, *prescriptionId, prescriptionDrug)
			if err != nil {
				return nil, apperror.InternalServerError(err)
			}
		}

		chat.Prescription.Id = prescriptionId
	}

	if file != nil {
		filePath, format, err := util.ValidateFile(*fileHeader, appconstant.ChatAttachmentUrl, []string{"png", "jpg", "jpeg", "pdf"}, 2000000)
		if err != nil {
			return nil, apperror.NewAppError(http.StatusBadRequest, err, err.Error())
		}

		AttachmentUrlUrl, err := util.UploadToCloudinary(file, *filePath)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		chat.Attachment.Format = format
		chat.Attachment.Url = &AttachmentUrlUrl
	}

	chatId, createdAt, err := chatRepo.PostOneChat(ctx, chat)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	chat.CreatedAt = &createdAt
	chat.Id = *chatId

	isListenerReady := u.findListenerBySender(entity.Participant{AccountId: chat.SenderAccountId, RoomId: chat.RoomId})
	if isListenerReady {
		u.chatChannel[chat.RoomId] <- chat
	}

	postMessageResponse := dto.ConvertToChatDTO(chat)

	return &postMessageResponse, nil
}

func (u *telemedicineUsecaseImpl) findListenerBySender(sender entity.Participant) bool {
	u.listenersLock.Lock()
	defer u.listenersLock.Unlock()

	for _, listener := range u.listeners {
		if sender.AccountId != listener.AccountId && sender.RoomId == listener.RoomId {
			return true
		}
	}

	return false
}

func (u *telemedicineUsecaseImpl) findListener(listener entity.Participant) bool {
	for _, l := range u.listeners {
		if l.AccountId == listener.AccountId && l.RoomId == listener.RoomId {
			return true
		}
	}

	return false
}

func (u *telemedicineUsecaseImpl) findListenerRoomIdByAccountId(listenerAccountId int64) *int64 {
	for _, l := range u.listeners {
		if l.AccountId == listenerAccountId {
			return &l.RoomId
		}
	}

	return nil
}

func (u *telemedicineUsecaseImpl) addListener(listener entity.Participant) {
	u.listenersLock.Lock()
	defer u.listenersLock.Unlock()

	u.listeners = append(u.listeners, listener)
}

func (u *telemedicineUsecaseImpl) removeListener(listener entity.Participant) {
	u.listenersLock.Lock()
	defer u.listenersLock.Unlock()

	for i, l := range u.listeners {
		if l.AccountId == listener.AccountId && l.RoomId == listener.RoomId {
			u.listeners = append(u.listeners[:i], u.listeners[i+1:]...)
			return
		}
	}
}

func (u *telemedicineUsecaseImpl) Listen(ctx context.Context, accountId int64, roomId int64) (*dto.Chat, error) {
	roomChat, err := u.chatRoomRepository.FindChatRoomById(ctx, roomId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if roomChat.DoctorAccountId != accountId && roomChat.UserAccountId != accountId {
		return nil, apperror.ForbiddenAction()
	}

	listener := entity.Participant{
		AccountId: accountId,
		RoomId:    roomChat.Id,
	}

	isAlreadyListen := u.findListener(listener)
	if isAlreadyListen {
		u.abortChannel <- listener
	}

	if u.chatChannel[listener.RoomId] == nil {
		u.chatChannel[listener.RoomId] = make(chan entity.Chat)
	}

	u.addListener(listener)

	resultChan := make(chan entity.Chat)
	errorChan := make(chan error)

	expiredAt := time.Now().Add(10 * time.Minute)

	if roomChat.ExpiredAt != nil {
		expiredAt = *roomChat.ExpiredAt
	}

	go func() {
		for {
			select {
			case chat := <-u.chatChannel[listener.RoomId]:
				if chat.Id == 0 {

					continue
				}

				if chat.SenderAccountId != listener.AccountId {
					resultChan <- chat

					return
				}

				u.chatChannel[chat.RoomId] <- chat

			case l := <-u.abortChannel:
				if listener.AccountId == l.AccountId && listener.RoomId == l.RoomId {
					errorChan <- apperror.AbortPreviousListenRequestError()

					return
				}

				u.abortChannel <- l

			case <-time.After(time.Until(expiredAt)):
				errorChan <- apperror.AbortPreviousListenRequestError()
				return
			}
		}
	}()

	select {
	case chatReceived := <-resultChan:
		chatResponse := dto.ConvertToChatDTO(chatReceived)
		u.removeListener(listener)

		return &chatResponse, nil
	case err := <-errorChan:
		u.removeListener(listener)

		return nil, err
	}
}

func (u *telemedicineUsecaseImpl) GetAllChat(ctx context.Context, accountId, roomId int64) (*dto.ChatRoom, error) {
	chatRoom, err := u.chatRoomRepository.FindChatRoomById(ctx, roomId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if chatRoom == nil {
		return nil, apperror.ChatRoomNotFoundError()
	}

	previousRoomId := u.findListenerRoomIdByAccountId(accountId)
	if previousRoomId != nil && *previousRoomId != roomId {
		u.abortChannel <- entity.Participant{RoomId: *previousRoomId, AccountId: accountId}
	}

	if chatRoom.DoctorAccountId != accountId && chatRoom.UserAccountId != accountId {
		return nil, apperror.ChatRoomNotFoundError()
	}

	doctorData, err := u.doctorRepository.FindDoctorByAccountId(ctx, chatRoom.DoctorAccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	chatRoom.DoctorCertificateUrl = doctorData.Certificate

	chats, err := u.chatRepository.GetAllChat(ctx, chatRoom.Id)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	var chatList []entity.Chat

	for _, chat := range chats {
		if chat.Prescription.Id != nil {
			prescriptionDrugList, err := u.prescriptionDrugRepository.GetAllPrescriptionDrug(ctx, *chat.Prescription.Id)
			if err != nil {
				return nil, err
			}

			chat.Prescription.PrescriptionDrugs = prescriptionDrugList
		}

		chatList = append(chatList, chat)
	}

	chatRoom.Chats = chatList

	chatRoomResponse := dto.ConvertToChatRoomDTO(*chatRoom)

	return &chatRoomResponse, nil
}

func (u *telemedicineUsecaseImpl) GetAllChatRoomPreview(ctx context.Context, accountId int64, role string) ([]dto.ChatRoomPreview, error) {
	chatRoomPreviewList, err := u.chatRoomRepository.GetAllChatRoomPreview(ctx, accountId, role)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	chatRoomPreviewResponse := dto.ConvertToChatRoomPreviewList(chatRoomPreviewList)

	return chatRoomPreviewResponse, nil
}

func (u *telemedicineUsecaseImpl) DoctorGetChatRequest(ctx context.Context, accountId int64) ([]dto.ChatRoomPreview, error) {
	chatRoomPreviewList, err := u.chatRoomRepository.DoctorGetChatRequest(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	chatRoomPreviewResponse := dto.ConvertToChatRoomPreviewList(chatRoomPreviewList)

	return chatRoomPreviewResponse, nil
}

func (u *telemedicineUsecaseImpl) SavePrescription(ctx context.Context, accountId, prescriptionId int64) error {
	prescription, err := u.prescriptionRepository.GetPrescriptionById(ctx, prescriptionId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if prescription == nil {
		return apperror.InvalidPrescriptionIdError()
	}

	if prescription.RedeemedAt != nil {
		return apperror.PrescriptionHasBeenRedeemedError()
	}

	if prescription.UserAccountId != accountId {
		return apperror.InvalidPrescriptionIdError()
	}

	err = u.prescriptionRepository.SetPrescriptionRedeemedNow(ctx, *prescription.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *telemedicineUsecaseImpl) GetAllPrescriptions(ctx context.Context, accountId int64, limit, page string) (*dto.PrescriptionResponseList, error) {
	limitInt, offsetInt, err := util.CheckPharmacyDrugPagination(page, limit)
	if err != nil {
		return nil, err
	}

	prescriptionList, err := u.prescriptionRepository.GetPrescriptionListByUserAccountId(ctx, accountId, limitInt, offsetInt)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	totalItem, err := u.prescriptionRepository.GetPrescriptionListByUserAccountIdTotalItem(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	totalPage := math.Ceil(float64(totalItem) / float64(limitInt))

	for i, prescription := range prescriptionList {
		prescriptionDrugs, err := u.prescriptionDrugRepository.GetAllPrescriptionDrug(ctx, *prescription.Id)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		prescriptionList[i].PrescriptionDrugs = prescriptionDrugs
	}

	response := dto.ConvertToPrescriptionResponseList(prescriptionList, totalItem, int(totalPage))

	return &response, nil
}

func (u *telemedicineUsecaseImpl) PrepareForCheckout(ctx context.Context, accountId, prescriptionId int64, addressIdString string) (*dto.PreapareForCheckoutResponse, error) {
	if prescriptionId < 1 {
		return nil, apperror.PrescriptionIdNotANumberError()
	}

	prescription, err := u.prescriptionRepository.GetPrescriptionById(ctx, prescriptionId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	addressId, err := strconv.Atoi(addressIdString)
	if err != nil {
		return nil, apperror.AddressIdInvalidError()
	}

	if prescription.OrderedAt != nil {
		return nil, apperror.PrescriptionHasBeenUsedError()
	}

	userCredential, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	userIdFromAddress, err := u.userAddressRepository.FindOneUserAddressById(ctx, int64(addressId))
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if userCredential.Id != *userIdFromAddress {
		return nil, apperror.ForbiddenAction()
	}

	userAddress, err := u.userAddressRepository.GetOneUserAddressByAddressId(ctx, int64(addressId))
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	userAddress.Id = int64(addressId)

	prescriptionDrugList, err := u.prescriptionDrugRepository.GetAllPrescriptionDrug(ctx, prescriptionId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	var checkoutPreparation entity.PrepareForCheckout

	checkoutPreparation.UserAddress = *userAddress

out:
	for _, prescriptionDrug := range prescriptionDrugList {
		if !prescriptionDrug.Drug.IsActive {
			return nil, apperror.DrugIsInactiveError()
		}

		pharmacy, drugQuantity, err := u.pharmacyDrugRepository.GetNearestAvailablePharmacyDrugByDrugId(ctx, prescriptionDrug.Drug.Id, userAddress.Id)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		if pharmacy == nil || drugQuantity == nil {
			return nil, apperror.NoDrugNearby()
		}

		for i, pharmacyDrug := range checkoutPreparation.Items {
			if pharmacyDrug.PharmacyId == pharmacy.Id {
				drugQuantity.Quantity = prescriptionDrug.Quantity
				checkoutPreparation.Items[i].DrugQuantities = append(checkoutPreparation.Items[i].DrugQuantities, *drugQuantity)
				checkoutPreparation.Items[i].Subtotal = checkoutPreparation.Items[i].Subtotal.Add(decimal.NewFromInt(int64(drugQuantity.Quantity)).Mul(drugQuantity.PharmacyDrug.Price))
				checkoutPreparation.Items[i].Weight = checkoutPreparation.Items[i].Weight.Add(decimal.NewFromInt(int64(prescriptionDrug.Quantity)).Mul(drugQuantity.PharmacyDrug.Drug.Weight))

				continue out
			}
		}

		weight := decimal.NewFromInt(int64(prescriptionDrug.Quantity)).Mul(drugQuantity.PharmacyDrug.Drug.Weight)

		avaiableCourierList, err := u.pharmacyRepository.GetAllCourierOptionsByPharmacyId(ctx, userAddress.Id, pharmacy.Id, weight.InexactFloat64())
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		drugQuantity.Quantity = prescriptionDrug.Quantity
		nearestPharmacyDrug := entity.PrepareForCheckoutItem{
			PharmacyId:      pharmacy.Id,
			PharmacyName:    pharmacy.Name,
			PharmacyAddress: pharmacy.Address,
			Distance:        pharmacy.Distance,
			DeliveryOptions: avaiableCourierList,
			Subtotal:        decimal.NewFromInt(int64(prescriptionDrug.Quantity)).Mul(drugQuantity.PharmacyDrug.Price),
			Weight:          weight,
			DrugQuantities:  []entity.DrugQuantity{*drugQuantity},
		}

		checkoutPreparation.Items = append(checkoutPreparation.Items, nearestPharmacyDrug)
	}

	response := dto.ConvertPrepareForCheckoutToResponse(checkoutPreparation.Items)

	return &response, nil
}

func (u *telemedicineUsecaseImpl) CheckoutFromPrescription(ctx context.Context, checkoutFromPrescriptionRequest dto.CheckoutFromPrescriptionRequest) (*int64, error) {
	user, err := u.userRepository.FindUserByAccountId(ctx, checkoutFromPrescriptionRequest.AccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if user == nil {
		return nil, apperror.UserNotFoundError()
	}

	prescription, err := u.prescriptionRepository.GetPrescriptionById(ctx, checkoutFromPrescriptionRequest.PrescriptionId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if prescription.OrderedAt != nil {
		return nil, apperror.PrescriptionHasBeenUsedError()
	}

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	orderRepo := tx.OrderRepository()
	cartRepo := tx.CartRepository()
	orderPharmacyRepo := tx.OrderPharmacyRepository()
	orderItemRepo := tx.OrderItemRepository()
	pharmacyDrugRepo := tx.PharmacyDrugRepo()
	stockChangeRepo := tx.StockChangeRepo()
	stockMutationRepo := tx.StockMutationRepo()
	prescriptionRepo := tx.PrescriptionRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	orderCheckoutRequest := dto.ConvertPrescriptionCheckoutRequest(checkoutFromPrescriptionRequest)

	for i, pharmacy := range checkoutFromPrescriptionRequest.Pharmacies {
		var cartItemIds []int64

		for _, phamacyDrugQuantity := range pharmacy.PharmacyDrugs {
			cartItemId, err := cartRepo.PostOneCart(ctx, checkoutFromPrescriptionRequest.AccountId, phamacyDrugQuantity.PharmacyDrugId, phamacyDrugQuantity.Quantity)
			if err != nil {
				return nil, apperror.InternalServerError(err)
			}

			cartItemIds = append(cartItemIds, *cartItemId)
		}

		orderCheckoutRequest.Pharmacies[i].CartItemIds = cartItemIds
	}

	orderId, err := orderRepo.PostOneOrder(ctx, user.Id, checkoutFromPrescriptionRequest.Address, checkoutFromPrescriptionRequest.TotalAmount)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	orderPharmacies, err := orderPharmacyRepo.PostOrderPharmacies(ctx, orderId, orderCheckoutRequest)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	for i, pharmacy := range orderCheckoutRequest.Pharmacies {
		orderPharmacies[i].CartItems, err = cartRepo.GetAllCartDetailByIds(ctx, pharmacy.CartItemIds)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
	}

	err = orderItemRepo.PostOrderItems(ctx, orderPharmacies)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	allCartItems := []entity.CartItemForCheckout{}
	for _, pharmacy := range orderPharmacies {
		allCartItems = append(allCartItems, pharmacy.CartItems...)
	}

	err = pharmacyDrugRepo.GetPharmacyDrugsByCartForUpdate(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	carts, err := cartRepo.GetAllCartsForChangesByCartIds(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	err = stockChangeRepo.PostStockChangesByCartIds(ctx, carts)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	pharmacyDrugs, err := pharmacyDrugRepo.UpdatePharmacyDrugsByCartId(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	alternatives, err := stockMutationRepo.GetPossibleStockMutation(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	stockMutationList := []entity.PossibleStockMutation{}
	stockChangesList := []entity.StockChange{}
	insufficientCartItems := []int64{}
	for _, pharmacyDrug := range pharmacyDrugs {
		if pharmacyDrug.Stock >= 0 {
			continue
		}
		stock := pharmacyDrug.Stock
		for _, alternative := range alternatives {
			if alternative.CartItemId != pharmacyDrug.CartId {
				continue
			}
			if stock+alternative.AlternativeStock < 0 {
				stock += alternative.AlternativeStock
				stockMutationList = append(stockMutationList, alternative)
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.OriginalPharmacyDrug,
					FinalStock: stock, Amount: alternative.AlternativeStock})
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.AlternativePharmacy,
					FinalStock: 0, Amount: -1 * alternative.AlternativeStock})
			} else {
				partialAlternative := alternative
				partialAlternative.AlternativeStock = stock * -1
				stockMutationList = append(stockMutationList, partialAlternative)
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.OriginalPharmacyDrug,
					FinalStock: 0, Amount: partialAlternative.AlternativeStock})
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.AlternativePharmacy,
					FinalStock: alternative.AlternativeStock - partialAlternative.AlternativeStock,
					Amount:     -1 * partialAlternative.AlternativeStock})
				stock = 0
				break
			}
		}
		if stock < 0 {
			insufficientCartItems = append(insufficientCartItems, pharmacyDrug.CartId)
		}
	}

	if len(insufficientCartItems) > 0 {
		err = apperror.InsufficientStockDuringCheckoutError(insufficientCartItems)
		return nil, err
	}

	if len(stockMutationList) > 0 {
		err = stockMutationRepo.PostStockMutations(ctx, stockMutationList)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		err = stockChangeRepo.PostStockChangesFromMutation(ctx, stockChangesList)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		err = pharmacyDrugRepo.UpdatePharmacyDrugsForStockMutation(ctx, stockChangesList)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
	}

	err = cartRepo.DeleteCarts(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	err = prescriptionRepo.SetPrescriptionOrderedAtNow(ctx, *prescription.Id)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &orderId, nil
}

func (u *telemedicineUsecaseImpl) CloseChatRoom(ctx context.Context, userAccountId, roomId int64) error {
	chatRoom, err := u.chatRoomRepository.FindChatRoomById(ctx, roomId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if chatRoom == nil {
		return apperror.ChatRoomNotFoundError()
	}

	if chatRoom.UserAccountId != userAccountId {
		return apperror.ForbiddenAction()
	}

	if chatRoom.ExpiredAt != nil && chatRoom.ExpiredAt.Before(time.Now()) {
		return apperror.ChatRoomAlreadyClosedError()
	}

	err = u.chatRoomRepository.CloseChatRoom(ctx, chatRoom.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	for _, l := range u.listeners {
		if l.RoomId == chatRoom.Id {
			u.abortChannel <- l
		}
	}

	return nil
}
