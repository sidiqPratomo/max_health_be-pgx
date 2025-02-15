package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sidiqPratomo/max-health-backend/config"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/handler"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func createRouter(log *logrus.Logger, config *config.Config) *gin.Engine {
	db := database.ConnectDB(config, log)

	accountRepository := repository.NewAccountRepositoryPostgres(db)
	cartRepository := repository.NewCartRepositoryPostgres(db)
	userRepository := repository.NewUserRepositoryPostgres(db)
	doctorRepository := repository.NewDoctorRepositoryPostgres(db)
	verificationCodeRepository := repository.NewVerificationCodeRepositoryPostgres(db)
	userAddressRepository := repository.NewUserAddressRepositoryPostgres(db)
	refreshTokenRepository := repository.NewRefreshTokenRepositoryPostgres(db)
	resetPasswordTokenRepository := repository.NewResetPasswordTokenRepositoryPostgres(db)
	pharmacyManagerRepository := repository.NewpharmacyManagerRepositoryPostgres(db)
	addressRepository := repository.NewAddressRepositoryPostgres(db)
	drugRepository := repository.NewDrugRepositoryPostgres(db)
	drugFormRepository := repository.NewDrugFormRepositoryPostgres(db)
	drugClassificationRepository := repository.NewDrugClassificationRepositoryPostgres(db)
	drugPharmacyRepository := repository.NewDrugPharmacyRepositoryPostgres(db)
	categoryRepository := repository.NewCategoryRepositoryPostgres(db)
	chatRepository := repository.NewChatRepositoryPostgres(db)
	chatRoomRepository := repository.NewChatRoomRepositoryPostgres(db)
	doctorSpecializationRepository := repository.NewDoctorSpecializationRepositoryPostgres(db)
	prescriptionRepository := repository.NewPrescriptionRepositoryPostgres(db)
	prescriptionDrugRepository := repository.NewPrescriptionDrugRepositoryPostgres(db)
	orderRepository := repository.NewOrderRepositoryPostgres(db)
	orderPharmacyRepository := repository.NewOrderPharmacyRepositoryPostgres(db)
	pharmacyRepository := repository.NewPharmacyRepositoryPostgres(db)
	courierRepository := repository.NewCourierRepositoryPostgres(db)
	orderItemRepository := repository.NewOrderItemRepositoryPostgres(db)
	stockRepository := repository.NewStockChangeRepositoryPostgres(db)
	transaction := repository.NewSqlTransaction(db)
	emailHelper := util.NewEmailHelperIpl(config)
	jwtAuthentication := util.JwtAuthentication{
		Config: *config,
		Method: jwt.SigningMethodHS256,
	}
	hashHelper := &util.HashHelperImpl{}

	authenticationUsecase := usecase.NewAuthenticationUsecaseImpl(usecase.AuthenticationUsecaseImplOpts{
		DrugRepository:               &drugRepository,
		AccountRepository:            &accountRepository,
		CartRepository:               &cartRepository,
		DoctorRepository:             &doctorRepository,
		UserRepository:               &userRepository,
		VerificationCodeRepository:   &verificationCodeRepository,
		RefreshTokenRepositoy:        &refreshTokenRepository,
		ResetPasswordTokenRepository: &resetPasswordTokenRepository,
		Transaction:                  transaction,
		HashHelper:                   hashHelper,
		JwtHelper:                    jwtAuthentication,
		EmailHelper:                  &emailHelper,
	})

	userUsecase := usecase.NewUserUsecaseImpl(&accountRepository, transaction, &userRepository, &userAddressRepository, &util.HashHelperImpl{})
	doctorUsecase := usecase.NewDoctorUsecaseImpl(&accountRepository, &doctorRepository, &doctorSpecializationRepository, transaction, &util.HashHelperImpl{})
	userAddressUsecase := usecase.NewUserAddressUsecaseImpl(&userRepository, &userAddressRepository, &addressRepository, transaction)
	partnerUsecase := usecase.NewPartnerUsecaseImpl(usecase.PartnerUsecaseImplOpts{
		AccountRepository:         &accountRepository,
		PharmacyManagerRepository: &pharmacyManagerRepository,
		Transaction:               transaction,
		HashHelper:                hashHelper,
		EmailHelper:               &emailHelper,
	})
	addressUsecase := usecase.NewAddressUsecaseImpl(&addressRepository)
	categoryUsecase := usecase.NewCategoryUsecaseImpl(&categoryRepository)

	drugUsecase := usecase.NewDrugUsecaseImpl(transaction, &drugRepository, &drugPharmacyRepository, &drugClassificationRepository, &drugFormRepository, &categoryRepository, &pharmacyRepository)
	drugFormUsecase := usecase.NewdrugFormUsecaseImpl(&drugFormRepository)
	drugClassificationUsecase := usecase.NewDrugClassificationUsecaseImpl(&drugClassificationRepository)
	telemedicineUsecase := usecase.NewTelemedicineUsecaseImpl(
		&chatRoomRepository,
		&chatRepository,
		&userRepository,
		&doctorRepository,
		&drugPharmacyRepository,
		&prescriptionDrugRepository,
		&prescriptionRepository,
		&cartRepository,
		&orderRepository,
		&userAddressRepository,
		&pharmacyRepository,
		transaction,
	)

	pharmacyUsecase := usecase.NewPharmacyUsecaseImpl(&pharmacyManagerRepository, &pharmacyRepository, &drugPharmacyRepository, &addressRepository, &courierRepository, &orderPharmacyRepository, transaction)

	cartUsecase := usecase.NewCartUsecaseImpl(&drugPharmacyRepository, &userRepository, &userAddressRepository, &cartRepository)
	orderUsecase := usecase.NewOrderUsecaseImpl(transaction, &userRepository, &orderRepository, &orderPharmacyRepository)
	orderPharmacyUsecase := usecase.NewOrderPharmacyUsecaseImpl(transaction, &orderPharmacyRepository, &orderItemRepository, &userRepository, &pharmacyManagerRepository)
	reportUsecase := usecase.NewreportUsecaseImpl(&orderItemRepository, &pharmacyRepository, &pharmacyManagerRepository)
	stockUsecase := usecase.NewStockUsecaseImpl(&stockRepository, &pharmacyManagerRepository)

	pingHandler := handler.NewPingHandler(handler.PingHandlerOpts{})
	authenticationHandler := handler.NewAuthenticationHandler(&authenticationUsecase)
	userHandler := handler.NewUserHandler(&userUsecase)
	doctorHandler := handler.NewDoctorHandler(&doctorUsecase)
	userAddressHandler := handler.NewUserAddressHandler(&userAddressUsecase)
	partnerHandler := handler.NewPartnerHandler(&partnerUsecase)
	addressHandler := handler.NewAddressHandler(&addressUsecase)
	drugHandler := handler.NewDrugHandler(&drugUsecase)
	cartHandler := handler.NewCartHandler(&cartUsecase)
	drugFormHandler := handler.NewDrugFormHandler(&drugFormUsecase)
	drugClassificationHandler := handler.NewDrugClassificationHandler(&drugClassificationUsecase)
	categoryHandler := handler.NewCategoryHandler(&categoryUsecase)
	telemedicineHandler := handler.NewTelemedicineHandler(&telemedicineUsecase)
	orderHandler := handler.NewOrderHandler(&orderUsecase)
	pharmacyHandler := handler.NewPharmacyHandler(&pharmacyUsecase)
	orderPharmacyHandler := handler.NewOrderPharmacyHandler(&orderPharmacyUsecase)
	reportHandler := handler.NewReportHandler(&reportUsecase)
	stockHandler := handler.NewStockHandler(&stockUsecase)

	return newRouter(
		routerOpts{
			Ping:               pingHandler,
			Authentication:     &authenticationHandler,
			User:               &userHandler,
			UserAddress:        &userAddressHandler,
			Doctor:             &doctorHandler,
			Partner:            &partnerHandler,
			Address:            &addressHandler,
			Drug:               &drugHandler,
			Category:           &categoryHandler,
			DrugForm:           &drugFormHandler,
			DrugClassification: &drugClassificationHandler,
			Cart:               &cartHandler,
			Telemedicine:       &telemedicineHandler,
			Order:              &orderHandler,
			Pharmacy:           &pharmacyHandler,
			OrderPharmacy:      &orderPharmacyHandler,
			Report:             &reportHandler,
			Stock:              &stockHandler,
		},
		utilOpts{
			JwtHelper: jwtAuthentication,
		},
		config,
		log,
	)
}

func Init() {
	log := util.NewLogger()

	config := config.Init(log)

	router := createRouter(log, config)

	srv := http.Server{
		Handler: router,
		Addr:    config.Port,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 10)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GracefulPeriod)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown: %s", err.Error())
	}

	<-ctx.Done()

	log.Infof("Timeout of " + strconv.Itoa(config.GracefulPeriod) + " seconds")
	log.Info("Server exiting")
}
