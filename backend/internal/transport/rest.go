package transport

import "github.com/Spiderpig02/AnkiStatsShower/internal/database"

type Handler struct {
	database      *database.Database
	SystsemStatus int
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{
		database: db,
	}
}
