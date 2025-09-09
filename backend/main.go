package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ---------------------------
// Models
// ---------------------------
type Entry struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID string `json:"user_id" gorm:"index"`
	Date   string `json:"date"`
	Number int    `json:"number"`
}

// ---------------------------
// Request types
// ---------------------------
type UploadRequest struct {
	UserID string       `json:"user_id"`
	Data   []DateNumber `json:"data"`
}

type DateNumber struct {
	Date   string `json:"date"`
	Number int    `json:"number"`
}

// ---------------------------
// Globals
// ---------------------------
var db *gorm.DB

// ---------------------------
// Handlers
// ---------------------------

// Upload a list of dates + numbers
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var req UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, dn := range req.Data {
		entry := Entry{
			UserID: req.UserID,
			Date:   dn.Date,
			Number: dn.Number,
		}
		db.Create(&entry)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Get data for a specific user
func getHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	var entries []Entry
	db.Where("user_id = ?", userID).Find(&entries)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// ---------------------------
// Main
// ---------------------------
func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Auto-migrate tables
	db.AutoMigrate(&Entry{})

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/get", getHandler)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
