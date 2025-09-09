package database

import (
	"log"
	"log/slog"

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

func (d *Database) GetUserByID(userID string) (*models.Entry, error) {
	var entry models.Entry
	result := d.client.Where("user_id = ?", userID).First(&entry)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entry, nil
}

func (d *Database) CreateUser(entry *models.Entry) error {
	result := d.client.Create(entry)
	if result.Error != nil {
		return result.Error
	}
	config.Logger.Info("Created new user", slog.String("user_id", entry.UserID))
	return nil
}

func (d *Database) PostData(pData *models.PostDataRequest) error {
	var entry models.Entry
	result := d.client.Where("secret_key = ?", pData.SecretKey).First(&entry)
	if result.Error != nil {
		log.Println("Error finding entry with secret key:", result.Error)
		return result.Error
	}

	entry.Data = pData.Data
	saveResult := d.client.Save(&entry)
	if saveResult.Error != nil {
		log.Println("Error updating entry data:", saveResult.Error)
		return saveResult.Error
	}

	config.Logger.Info("Updated data for user", slog.String("user_id", entry.UserID))
	return nil
}
