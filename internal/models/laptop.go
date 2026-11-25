package models

import "time"

// Laptop representa una laptop disponible para préstamo.
type Laptop struct {
	ID             int       `db:"id" json:"id"`
	Brand          string    `db:"brand" json:"brand"`
	Model          string    `db:"model" json:"model"`
	SerialNumber   string    `db:"serial_number" json:"serial_number"`
	TotalUnits     int       `db:"total_units" json:"total_units"`
	AvailableUnits int       `db:"available_units" json:"available_units"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// LaptopRequest es el payload para crear una laptop.
type LaptopRequest struct {
	Brand        string `json:"brand" binding:"required" example:"Dell"`
	Model        string `json:"model" binding:"required" example:"Latitude 5420"`
	SerialNumber string `json:"serial_number" binding:"required" example:"DELL-001"`
	TotalUnits   int    `json:"total_units" binding:"required,min=1" example:"5"`
}
