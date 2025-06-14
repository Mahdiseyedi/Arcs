package db

import (
	"arcs/internal/configs"
	"arcs/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg configs.Config) (DB *Database) {
	var loger logger.Interface

	if cfg.Basic.Environment == "dev" {
		loger = logger.Default.LogMode(logger.Info)
	} else {
		loger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.ConnectionString), &gorm.Config{
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
		Logger:                 loger,
	})
	if err != nil {
		log.Fatalf("[DATABASE] Failed to connect database: [%v]", err)
	}

	database := &Database{
		DB: db,
	}
	
	database.autoMigrate()

	return database
}

func (db *Database) autoMigrate() {
	if err := db.DB.AutoMigrate(
		&models.User{},
		&models.Order{},
		&models.SMS{},
	); err != nil {
		log.Fatalf("[DATABASE] Failed to run auto migrate: [%v]", err)
	}
}
