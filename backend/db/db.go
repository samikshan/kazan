package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	"github.com/samikshan/kazan/backend"
	"github.com/samikshan/kazan/backend/models"
)

func New() *gorm.DB {
	cfg := backend.Cfg
	// Connect to database
	dbConnArgs := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	db, err := gorm.Open("postgres", dbConnArgs)
	if err != nil {
		log.WithError(err).Fatal("db connection failed")
	}

	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return db
}

//TODO: err check
func AutoMigrate(db *gorm.DB) {
	log.Info("migrate shit")
	db.AutoMigrate(
		&models.User{},
	)
}
