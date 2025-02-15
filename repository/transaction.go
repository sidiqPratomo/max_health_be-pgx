package repository

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Transaction interface {
	Rollback() error
	Commit() error
	BeginTx() (Transaction, error)
	AccountRepository() AccountRepository
	UserRepository() UserRepository
	DoctorRepository() DoctorRepository
	VerificationCodeRepository() VerificationCodeRepository
	RefreshTokenRepository() RefreshTokenRepository
	UserAddressRepository() UserAddressRepository
	ResetPasswordTokenRepository() ResetPasswordTokenRepository
	PharmacyManagerRepository() PharmacyManagerRepository
	AddressRepository() AddressRepository
	PrescriptionRepository() PrescriptionRepository
	PrescriptionDrugRepository() PrescriptionDrugRepository
	ChatRepository() ChatRepository
	OrderRepository() OrderRepository
	OrderPharmacyRepository() OrderPharmacyRepository
	OrderItemRepository() OrderItemRepository
	CartRepository() CartRepository
	PharmacyDrugRepo() PharmacyDrugRepository
	StockChangeRepo() StockChangeRepository
	StockMutationRepo() StockMutationRepository
	PharmacyRepository() PharmacyRepository
	PharmacyOperationalRepository() PharmacyOperationalRepository
	PharmacyCourierRepository() PharmacyCourierRepository
}

type SqlTransaction struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSqlTransaction(db *sql.DB) *SqlTransaction {
	return &SqlTransaction{
		db: db,
	}
}

func (s *SqlTransaction) BeginTx() (Transaction, error) {
	tx, err := s.db.Begin()
	return &SqlTransaction{db: s.db, tx: tx}, err
}

func (s *SqlTransaction) Rollback() error {
	return s.tx.Rollback()
}

func (s *SqlTransaction) Commit() error {
	return s.tx.Commit()
}

func (s *SqlTransaction) AccountRepository() AccountRepository {
	return &accountRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) UserRepository() UserRepository {
	return &userRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) DoctorRepository() DoctorRepository {
	return &doctorRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) VerificationCodeRepository() VerificationCodeRepository {
	return &verificationCodeRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) RefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) UserAddressRepository() UserAddressRepository {
	return &userAddressRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) ResetPasswordTokenRepository() ResetPasswordTokenRepository {
	return &resetPasswordTokenRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) PharmacyManagerRepository() PharmacyManagerRepository {
	return &pharmacyManagerRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) AddressRepository() AddressRepository {
	return &addressRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) PrescriptionRepository() PrescriptionRepository {
	return &prescriptionRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) OrderRepository() OrderRepository {
	return &orderRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) PrescriptionDrugRepository() PrescriptionDrugRepository {
	return &prescriptionDrugRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) OrderPharmacyRepository() OrderPharmacyRepository {
	return &orderPharmacyRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) ChatRepository() ChatRepository {
	return &chatRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) OrderItemRepository() OrderItemRepository {
	return &orderItemRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) CartRepository() CartRepository {
	return &cartRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) PharmacyDrugRepo() PharmacyDrugRepository {
	return &pharmacyDrugRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) StockChangeRepo() StockChangeRepository {
	return &stockChangeRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) StockMutationRepo() StockMutationRepository {
	return &stockMutationRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) PharmacyRepository() PharmacyRepository {
	return &pharmacyRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) PharmacyOperationalRepository() PharmacyOperationalRepository {
	return &pharmacyOperationalRepositoryPostgres{
		db: s.tx,
	}
}

func (s *SqlTransaction) PharmacyCourierRepository() PharmacyCourierRepository {
	return &pharmacyCourierRepositoryPostgres{
		db: s.tx,
	}
}
