package models

import "time"

type Entry struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	CreationTime time.Time `gorm:"autoCreateTime"`
	UserID       string    `json:"user_id"`
	SecretKey    string    `json:"secret_key"`
	Data         string    `json:"data"`
}
