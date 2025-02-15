package database

import (
	"context"
	// "database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sidiqPratomo/max-health-backend/config"
	"github.com/sirupsen/logrus"
)

// func ConnectDB(config *config.Config, log *logrus.Logger) *sql.DB {
// 	db, err := sql.Open("pgx", config.DbUrl)
// 	if err != nil {
// 		log.WithFields(logrus.Fields{
// 			"error": err.Error(),
// 		}).Fatal("error connecting to DB")
// 		return nil
// 	}

// 	if err = db.Ping(); err != nil {
// 		log.WithFields(logrus.Fields{
// 			"error": err.Error(),
// 		}).Fatal("error connecting to DB")
// 		return nil
// 	}

// 	return db
// }


func ConnectDB(config *config.Config, log *logrus.Logger) *pgxpool.Pool {
    dbpool, err := pgxpool.New(context.Background(), config.DbUrl)
    if err != nil {
        log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Fatal("error connecting to DB")
        return nil
    }

    if err = dbpool.Ping(context.Background()); err != nil {
        log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Fatal("error connecting to DB")
        return nil
    }

    return dbpool
}
