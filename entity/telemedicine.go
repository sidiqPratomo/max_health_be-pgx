package entity

import (
	"time"
)

type Chat struct {
	Id              int64
	RoomId          int64
	SenderAccountId int64
	Message         *string
	Attachment      Attachment
	Prescription    Prescription
	CreatedAt       *string
}

type Attachment struct {
	Url    *string
	Format *string
}

type ChatRoom struct {
	Id                   int64
	DoctorAccountId      int64
	UserAccountId        int64
	DoctorCertificateUrl string
	ExpiredAt            *time.Time
	ExpiredAtString      *string
	Chats                []Chat
}

type ChatRoomPreview struct {
	Id                    int64
	ParticipantName       string
	ParticipantPictureUrl string
	ExpiredAt             *time.Time
	LastChat              Chat
}

type Participant struct {
	AccountId int64
	RoomId    int64
}
