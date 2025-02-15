package database

import (
	"database/sql"

	"github.com/sidiqPratomo/max-health-backend/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func ConnectDB(config *config.Config, log *logrus.Logger) *sql.DB {
	db, err := sql.Open("pgx", config.DbUrl)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error connecting to DB")
		return nil
	}

	if err = db.Ping(); err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error connecting to DB")
		return nil
	}

	return db
}
