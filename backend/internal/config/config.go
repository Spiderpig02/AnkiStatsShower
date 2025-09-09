package config

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger
var ServerPort = "8080"
var DatabaseName = "data.db"

func Init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

}
