package transport

import (
	"log/slog"

	"github.com/Spiderpig02/AnkiStatsShower/internal/config"
	"github.com/Spiderpig02/AnkiStatsShower/internal/database"
	"github.com/Spiderpig02/AnkiStatsShower/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	database     *database.Database
	SystemStatus int
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{
		database: db,
	}
}

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			config.Logger.Warn("UserID is missing in the request")
			c.JSON(400, gin.H{"error": "UserID is required"})
			return
		}

		// Check if user already exists
		existingUser, err := h.database.GetUserByID(userID)
		if err == nil && existingUser != nil {
			config.Logger.Info("User already exists:", userID)
			c.JSON(200, gin.H{"message": "User already exists"})
			return
		}

		// Create new user
		secretKey := uuid.New().String() + ":" + userID
		newUser := &models.Entry{
			ID:        uuid.New().String(),
			UserID:    userID,
			SecretKey: secretKey,
		}

		err = h.database.CreateUser(newUser)
		if err != nil {
			config.Logger.Error("Error creating user:", err)
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}

		config.Logger.Info("Created new user:", userID)
		c.JSON(201, gin.H{
			"message":    "User created successfully",
			"user_id":    userID,
			"secret_key": newUser.SecretKey,
		})
	}
}

func (h *Handler) PostData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pData models.PostDataRequest
		if err := c.ShouldBindJSON(&pData); err != nil {
			config.Logger.Warn("Invalid request body:", slog.String("error", err.Error()))
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		err := h.database.PostData(&pData)
		if err != nil {
			config.Logger.Error("Error posting data:", err)
			c.JSON(500, gin.H{"error": "Failed to post data"})
			return
		}

		config.Logger.Info("Data posted successfully for secret key:", pData.SecretKey)
		c.JSON(200, gin.H{"message": "Data posted successfully"})
	}
}

func (h *Handler) GetUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			config.Logger.Warn("UserID is missing in the request")
			c.JSON(400, gin.H{"error": "UserID is required"})
			return
		}

		entry, err := h.database.GetUserByID(userID)
		if err != nil {
			config.Logger.Error("Error retrieving user data:", err)
			c.JSON(500, gin.H{"error": "Failed to retrieve user data"})
			return
		}
		if entry == nil {
			config.Logger.Info("User not found:", userID)
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		c.JSON(200, gin.H{
			"user_id": entry.UserID,
			"data":    entry.Data,
		})
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
