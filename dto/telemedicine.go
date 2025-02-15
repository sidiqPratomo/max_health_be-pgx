package dto

import (
	"time"

	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/shopspring/decimal"
)

type UserCreateRoomRequest struct {
	DoctorAccountId int64 `json:"doctor_account_id" binding:"required,gte=1"`
}

type DoctorJoinRoomRequest struct {
	RoomId int64 `json:"room_id" binding:"required,gte=1"`
}

type PostOneMessageRequest struct {
	RoomId            int64                     `json:"room_id" validate:"required,gte=1"`
	Message           string                    `json:"message" validate:"required,min=1"`
	PrescriptionDrugs []PrescriptionDrugRequest `json:"prescription_drugs"`
}

type UserCreateRoomResponse struct {
	RoomId *int64 `json:"room_id"`
}

type Chat struct {
	Id              int64                `json:"id"`
	RoomId          int64                `json:"room_id,omitempty"`
	SenderAccountId int64                `json:"sender_account_id,omitempty"`
	Message         *string              `json:"message"`
	Attachment      Attachment           `json:"attachment,omitempty"`
	Prescription    PrescriptionResponse `json:"prescription"`
	CreatedAt       *string              `json:"created_at"`
}

type Attachment struct {
	Url    *string `json:"url"`
	Format *string `json:"format"`
}

type ChatRoom struct {
	Id                   int64      `json:"id"`
	DoctorAccountId      int64      `json:"doctor_account_id"`
	UserAccountId        int64      `json:"user_account_id"`
	DoctorCertificateUrl string     `json:"doctor_certificate_url"`
	ExpiredAt            *time.Time `json:"expired_at,omitempty"`
	Chats                []Chat     `json:"chats"`
}

type ChatRoomPreview struct {
	Id                    int64      `json:"id"`
	ParticipantName       string     `json:"participant_name"`
	ParticipantPictureUrl string     `json:"participant_picture_url"`
	ExpiredAt             *time.Time `json:"expired_at,omitempty"`
	LastChat              Chat       `json:"last_chat"`
}

type CheckoutFromPrescriptionRequest struct {
	AccountId      int64
	PrescriptionId int64                                     `json:"prescription_id" binding:"required,gte=1"`
	Address        string                                    `json:"address" binding:"required"`
	TotalAmount    int                                       `json:"total_amount" binding:"required"`
	Pharmacies     []PharmacyCheckoutFromPrescriptionRequest `json:"pharmacies" binding:"required,min=1"`
}

type PreapareForCheckoutResponse struct {
	PharmacyDrugs []PreapareForCheckoutItemResponse `json:"pharmacy_drugs"`
}

type PreapareForCheckoutItemResponse struct {
	PharmacyId      int64                     `json:"pharmacy_id"`
	PharmacyName    string                    `json:"pharmacy_name"`
	PharmacyAddress string                    `json:"pharmacy_address"`
	Distance        float64                   `json:"distance"`
	Subtotal        decimal.Decimal           `json:"subtotal"`
	Couriers        []entity.AvailableCourier `json:"couriers"`
	DrugQuantities  []DrugQuantity            `json:"drug_quantities"`
}

type DrugQuantity struct {
	PharmacyDrug DetailPharmacyDrug `json:"pharmacy_drug"`
	Quantity     int                `json:"quantity"`
}

func ConvertToChatDTO(chat entity.Chat) Chat {
	return Chat{
		Id:              chat.Id,
		RoomId:          chat.RoomId,
		SenderAccountId: chat.SenderAccountId,
		Message:         chat.Message,
		Attachment:      (Attachment)(chat.Attachment),
		Prescription:    ConvertToPrescriptionResponse(chat.Prescription),
		CreatedAt:       chat.CreatedAt,
	}
}

func ConvertToChatListDTO(chatList []entity.Chat) []Chat {
	var chatListDTO []Chat

	for _, chat := range chatList {
		chatListDTO = append(chatListDTO, ConvertToChatDTO(chat))
	}

	return chatListDTO
}

func ConvertToChatRoomDTO(chatRoom entity.ChatRoom) ChatRoom {
	return ChatRoom{
		Id:                   chatRoom.Id,
		DoctorAccountId:      chatRoom.DoctorAccountId,
		UserAccountId:        chatRoom.UserAccountId,
		DoctorCertificateUrl: chatRoom.DoctorCertificateUrl,
		ExpiredAt:            chatRoom.ExpiredAt,
		Chats:                ConvertToChatListDTO(chatRoom.Chats),
	}
}

func ConvertToChatRoomPreview(chatRoomPreview entity.ChatRoomPreview) ChatRoomPreview {
	return ChatRoomPreview{
		Id:                    chatRoomPreview.Id,
		ParticipantName:       chatRoomPreview.ParticipantName,
		ParticipantPictureUrl: chatRoomPreview.ParticipantPictureUrl,
		ExpiredAt:             chatRoomPreview.ExpiredAt,
		LastChat:              ConvertToChatDTO(chatRoomPreview.LastChat),
	}
}

func ConvertToChatRoomPreviewList(chatRoomPreviewList []entity.ChatRoomPreview) []ChatRoomPreview {
	var chatRoomPreviewListDTO []ChatRoomPreview

	for _, chatRoomPreview := range chatRoomPreviewList {
		chatRoomPreviewListDTO = append(chatRoomPreviewListDTO, ConvertToChatRoomPreview(chatRoomPreview))
	}

	return chatRoomPreviewListDTO
}

func ConvertToPrescriptionEntity(prescriptionDrugs []PrescriptionDrugRequest) entity.Prescription {
	var prescription entity.Prescription

	for _, prescriptionDrug := range prescriptionDrugs {
		prescription.PrescriptionDrugs = append(prescription.PrescriptionDrugs, entity.PrescriptionDrug{
			Id:       prescriptionDrug.Id,
			Drug:     ConvertPrescriptionDrugRequestToDrug(prescriptionDrug.Drug),
			Quantity: prescriptionDrug.Quantity,
			Note:     prescriptionDrug.Note,
		})
	}

	return prescription
}

func ConvertPostMessageRequestToChatEntity(dto PostOneMessageRequest) entity.Chat {
	return entity.Chat{
		RoomId:       dto.RoomId,
		Message:      &dto.Message,
		Prescription: ConvertToPrescriptionEntity(dto.PrescriptionDrugs),
	}
}

func ConvertToDrugQuantityDTO(drugQuantity entity.DrugQuantity) DrugQuantity {
	return DrugQuantity{
		PharmacyDrug: DetailPharmacyDrug{
			Id:    drugQuantity.PharmacyDrug.Id,
			Drug:  ConvertToDrugResponse(drugQuantity.PharmacyDrug.Drug),
			Price: drugQuantity.PharmacyDrug.Price,
		},
		Quantity: drugQuantity.Quantity,
	}
}

func ConverToDrugQuantityListDTO(drugQuantityList []entity.DrugQuantity) []DrugQuantity {
	var response []DrugQuantity

	for _, drugQuantity := range drugQuantityList {
		response = append(response, ConvertToDrugQuantityDTO(drugQuantity))
	}

	return response
}

func ConvertPreapareForCheckoutItemToResponse(checkoutItem entity.PrepareForCheckoutItem) PreapareForCheckoutItemResponse {
	return PreapareForCheckoutItemResponse{
		PharmacyId:      checkoutItem.PharmacyId,
		PharmacyName:    checkoutItem.PharmacyName,
		PharmacyAddress: checkoutItem.PharmacyAddress,
		Distance:        checkoutItem.Distance,
		Subtotal:        checkoutItem.Subtotal,
		Couriers:        checkoutItem.DeliveryOptions,
		DrugQuantities:  ConverToDrugQuantityListDTO(checkoutItem.DrugQuantities),
	}
}

func ConvertPrepareForCheckoutToResponse(checkoutItemList []entity.PrepareForCheckoutItem) PreapareForCheckoutResponse {
	var response PreapareForCheckoutResponse

	for _, checkoutItem := range checkoutItemList {
		response.PharmacyDrugs = append(response.PharmacyDrugs, ConvertPreapareForCheckoutItemToResponse(checkoutItem))
	}

	return response
}
