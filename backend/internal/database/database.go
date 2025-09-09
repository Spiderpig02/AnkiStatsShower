package database

import (
	"log"

	"github.com/Spiderpig02/AnkiStatsShower/internal/config"
	"github.com/Spiderpig02/AnkiStatsShower/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	client *gorm.DB
}

func Connect2Database(name string) *Database {
	client, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Auto-migrate tables
	client.AutoMigrate(&models.Entry{})
	return &Database{client}
}

func (d *Database) Disconnect() error {
	sqlDB, err := d.client.DB()
	if err != nil {
		return err
	}
	config.Logger.Info("Disconnected from the database")
	return sqlDB.Close()
}
