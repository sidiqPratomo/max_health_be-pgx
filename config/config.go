package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port                string
	FEPort              string
	DbUrl               string
	Issuer              string
	SendEmailIdentity   string
	SendEmailUsername   string
	SendEmailPassword   string
	SendEmailHost       string
	SendEmailPort       string
	VerifSecret         string
	AccessSecret        string
	RefreshSecret       string
	ResetPasswordSecret string
	RajaOngkirApiKey    string
	HashCost            int
	GracefulPeriod      int
}

func Init(log *logrus.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	hashCost, err := strconv.Atoi(os.Getenv("HASH_COST"))
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": "HASH_COST must be integer",
		}).Fatal("error loading .env file")
	}

	gracefulPeriod, err := strconv.Atoi(os.Getenv("GRACEFUL_PERIOD"))
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": "GRACEFUL_PERIOD must be integer",
		}).Fatal("error loading .env file")
	}

	return &Config{
		Port:                os.Getenv("BE_PORT"),
		FEPort:              os.Getenv("FE_PORT"),
		DbUrl:               os.Getenv("DATABASE_URL"),
		Issuer:              os.Getenv("ISSUER"),
		SendEmailIdentity:   os.Getenv("SEND_EMAIL_IDENTITY"),
		SendEmailUsername:   os.Getenv("SEND_EMAIL_USERNAME"),
		SendEmailPassword:   os.Getenv("SEND_EMAIL_PASSWORD"),
		SendEmailHost:       os.Getenv("SEND_EMAIL_HOST"),
		SendEmailPort:       os.Getenv("SEND_EMAIL_PORT"),
		VerifSecret:         os.Getenv("VERIFICATION_CODE_SECRET_KEY"),
		AccessSecret:        os.Getenv("ACCESS_TOKEN_SECRET_KEY"),
		RefreshSecret:       os.Getenv("REFRESH_TOKEN_SECRET_KEY"),
		ResetPasswordSecret: os.Getenv("RESET_PASSWORD_SECRET_KEY"),
		RajaOngkirApiKey:    os.Getenv("RAJA_ONGKIR_API_KEY"),
		HashCost:            hashCost,
		GracefulPeriod:      gracefulPeriod,
	}
}
