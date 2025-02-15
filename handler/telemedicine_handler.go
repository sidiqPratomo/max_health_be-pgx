package handler

import (
	"encoding/json"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TelemedicineHandler struct {
	telemedicineUsecase usecase.TelemedicineUsecase
}

func NewTelemedicineHandler(telemedicineUsecase usecase.TelemedicineUsecase) TelemedicineHandler {
	return TelemedicineHandler{
		telemedicineUsecase: telemedicineUsecase,
	}
}

func (h *TelemedicineHandler) UserCreateRoom(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	userAccountId, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	var createRoomRequest dto.UserCreateRoomRequest

	err := ctx.ShouldBindJSON(&createRoomRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	roomId, err := h.telemedicineUsecase.UserCreateRoom(ctx.Request.Context(), userAccountId.(int64), createRoomRequest.DoctorAccountId)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, dto.UserCreateRoomResponse{RoomId: roomId})
}

func (h *TelemedicineHandler) DoctorJoinRoom(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	doctorAccountId, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	var joinRoomRequest dto.DoctorJoinRoomRequest

	err := ctx.ShouldBindJSON(&joinRoomRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.telemedicineUsecase.DoctorJoinRoom(ctx.Request.Context(), doctorAccountId.(int64), joinRoomRequest.RoomId)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *TelemedicineHandler) PostOneMessage(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var postOneMessageRequest dto.PostOneMessageRequest

	accountId, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	role, exists := ctx.Get(appconstant.Role)
	if !exists || (role != appconstant.UserRoleName && role != appconstant.DoctorRoleName) {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	value := ctx.Request.FormValue("data")
	if value == "" {
		ctx.Error(apperror.EmptyDataError())
		return
	}

	err := json.Unmarshal([]byte(value), &postOneMessageRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(postOneMessageRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	if len(postOneMessageRequest.PrescriptionDrugs) > 0 {
		for _, prescriptionDrug := range postOneMessageRequest.PrescriptionDrugs {
			err = validator.New().Struct(prescriptionDrug)
			if err != nil {
				ctx.Error(err)
				return
			}

			err = validator.New().Struct(prescriptionDrug.Drug)
			if err != nil {
				ctx.Error(err)
				return
			}
		}
	}

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		if file != nil {
			ctx.Error(err)
			return
		}
	}

	response, err := h.telemedicineUsecase.PostOneMessage(ctx, accountId.(int64), postOneMessageRequest, file, fileHeader)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, response)
}

func (h *TelemedicineHandler) Listen(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	role, exists := ctx.Get(appconstant.Role)
	if !exists || (role != appconstant.UserRoleName && role != appconstant.DoctorRoleName) {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	roomIdStr := ctx.Param(appconstant.RoomIdString)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	chat, err := h.telemedicineUsecase.Listen(ctx.Request.Context(), accountId.(int64), int64(roomId))
	if err != nil {
		if err == apperror.AbortPreviousListenRequestError() {
			util.ResponseOK(ctx, nil)
			return
		}

		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, chat)
}

func (h *TelemedicineHandler) GetAllChat(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	role, exists := ctx.Get(appconstant.Role)
	if !exists || (role != appconstant.UserRoleName && role != appconstant.DoctorRoleName) {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	roomIdString := ctx.Param(appconstant.RoomIdString)
	roomId, err := strconv.Atoi(roomIdString)
	if err != nil {
		ctx.Error(err)
		return
	}

	chatRoom, err := h.telemedicineUsecase.GetAllChat(ctx.Request.Context(), accountId.(int64), int64(roomId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, chatRoom)
}

func (h *TelemedicineHandler) GetAllChatRoomPreview(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	role, exists := ctx.Get(appconstant.Role)
	if !exists || (role != appconstant.UserRoleName && role != appconstant.DoctorRoleName) {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	chatRoomList, err := h.telemedicineUsecase.GetAllChatRoomPreview(ctx.Request.Context(), accountId.(int64), role.(string))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, chatRoomList)
}

func (h *TelemedicineHandler) DoctorGetChatRequest(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	chatRoomList, err := h.telemedicineUsecase.DoctorGetChatRequest(ctx.Request.Context(), accountId.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, chatRoomList)
}

func (h *TelemedicineHandler) SavePrescription(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	prescriptionIdString := ctx.Param(appconstant.PrescriptionIdString)

	prescriptionId, err := strconv.Atoi(prescriptionIdString)
	if err != nil {
		ctx.Error(apperror.PrescriptionIdNotANumberError())
		return
	}

	err = h.telemedicineUsecase.SavePrescription(ctx.Request.Context(), accountId.(int64), int64(prescriptionId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *TelemedicineHandler) GetAllPrescriptions(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	limit := ctx.Query(appconstant.Limit)
	page := ctx.Query(appconstant.Page)

	prescriptionList, err := h.telemedicineUsecase.GetAllPrescriptions(ctx.Request.Context(), accountId.(int64), limit, page)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, prescriptionList)
}

func (h *TelemedicineHandler) PreapereForCheckout(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	prescriptionIdString := ctx.Param(appconstant.PrescriptionIdString)
	addressId := ctx.Query(appconstant.AddressIdString)

	prescriptionId, err := strconv.Atoi(prescriptionIdString)
	if err != nil {
		ctx.Error(apperror.PrescriptionIdNotANumberError())
		return
	}

	nearestPharmacyDrugList, err := h.telemedicineUsecase.PrepareForCheckout(ctx.Request.Context(), accountId.(int64), int64(prescriptionId), addressId)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, *nearestPharmacyDrugList)
}

func (h *TelemedicineHandler) CheckoutFromPrescription(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	var checkoutFromPrescriptionRequest dto.CheckoutFromPrescriptionRequest

	err := ctx.ShouldBindJSON(&checkoutFromPrescriptionRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	for _, request := range checkoutFromPrescriptionRequest.Pharmacies {
		err = validator.New().Struct(request)
		if err != nil {
			ctx.Error(err)
			return
		}

		for _, pharmacyDrugQuantity := range request.PharmacyDrugs {
			err = validator.New().Struct(pharmacyDrugQuantity)
			if err != nil {
				ctx.Error(err)
				return
			}
		}
	}

	checkoutFromPrescriptionRequest.AccountId = accountId.(int64)

	orderId, err := h.telemedicineUsecase.CheckoutFromPrescription(ctx.Request.Context(), checkoutFromPrescriptionRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, dto.OrderCheckoutResponse{OrderId: *orderId})
}

func (h *TelemedicineHandler) CloseChatRoom(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	roomIdStr := ctx.Param(appconstant.RoomIdString)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	err = h.telemedicineUsecase.CloseChatRoom(ctx.Request.Context(), accountId.(int64), int64(roomId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}
