package transport

import (
	"github.com/Spiderpig02/AnkiStatsShower/internal/config"
	"github.com/Spiderpig02/AnkiStatsShower/internal/database"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	database      *database.Database
	SystsemStatus int
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{
		database: db,
	}
}

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *Handler) DoesUserExist(userId string) bool {
	_, err := h.database.GetUserByID(userId)
	if err != nil {
		config.Logger.Error("Error checking user existence:", err)
		return false
	}
	config.Logger.Info("User exists:", userId)
	return true
}
