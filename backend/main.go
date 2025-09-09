package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Spiderpig02/AnkiStatsShower/internal/config"
	"github.com/Spiderpig02/AnkiStatsShower/internal/database"
	"github.com/Spiderpig02/AnkiStatsShower/internal/transport"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config.Init()
	r := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	r.Use(
		cors.New(corsConfig),
		gin.LoggerWithWriter(gin.DefaultWriter, "/status"),
		gin.Recovery(),
	)

	database := database.Connect2Database("data.db")
	restHandler := transport.NewHandler(database)
	setupRoutes(r, restHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := r.Run(":" + config.ServerPort); err != nil {
			config.Logger.Warn("Failed to run server:", slog.String("error", err.Error()))
			stop <- syscall.SIGTERM
		}
	}()

	<-stop
	config.Logger.Info("Shutting down the server...")
	database.Disconnect()
	config.Logger.Info("Server gracefully stopped")
}

func setupRoutes(r *gin.Engine, handler *transport.Handler) {
	// REST routes
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Successfully started and running the AnkiStatsShower backend\n")
	})
}
