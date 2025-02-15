package entity

import "time"

type PrescriptionDrug struct {
	Id       int64
	Drug     Drug
	Quantity int
	Note     string
}

type Prescription struct {
	Id                *int64
	UserAccountId     int64
	UserName          string
	DoctorAccountId   int64
	DoctorName        string
	PrescriptionDrugs []PrescriptionDrug
	RedeemedAt        *time.Time
	OrderedAt         *time.Time
	CreatedAt         *time.Time
}
