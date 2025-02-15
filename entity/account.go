package entity

import "time"

type Role struct {
	Id   int64
	Name string
}

type Account struct {
	Id             int64
	Email          string
	Password       string
	RoleId         int64
	RoleName       string
	Name           string
	ProfilePicture string
	VerifiedAt     *time.Time
}

type VerificationCode struct {
	Id        int64
	AccountId int64
	Code      string
	ExpiredAt time.Time
}

type ResetPasswordToken struct {
	Id        int64
	AccountId int64
	Token     string
	ExpiredAt time.Time
}

type RefreshToken struct {
	Id        int64
	AccountId int64
	Token     string
	ExpiredAt time.Time
}
