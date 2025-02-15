package entity

import (
	"github.com/shopspring/decimal"
)

type Doctor struct {
	Id                 int64
	AccountId          int64
	Certificate        string
	FeePerPatient      decimal.Decimal
	IsOnline           bool
	Experience         int
	SpecializationId   int64
	SpecializationName string
}

type DetailedDoctor struct {
	Id                 int64
	Email              string
	Name               string
	ProfilePicture     string
	Password           string
	FeePerPatient      decimal.Decimal
	Experience         int
	SpecializationId   int64
	SpecializationName string
}

type DoctorSpecialization struct {
	Id   int64
	Name string
}
