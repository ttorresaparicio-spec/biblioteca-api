package models

import "time"

// LaptopLoan representa el préstamo de una laptop a un usuario.
type LaptopLoan struct {
	ID         int        `db:"id" json:"id"`
	UserID     int        `db:"user_id" json:"user_id"`
	LaptopID   int        `db:"laptop_id" json:"laptop_id"`
	LoanedAt   time.Time  `db:"loaned_at" json:"loaned_at"`
	DueAt      time.Time  `db:"due_at" json:"due_at"`
	ReturnedAt *time.Time `db:"returned_at" json:"returned_at,omitempty"`
	Status     string     `db:"status" json:"status"`
}

// LaptopLoanRequest es el payload para solicitar una laptop.
type LaptopLoanRequest struct {
	LaptopID int `json:"laptop_id" binding:"required" example:"1"`
}

// LaptopLoanDetail combina el préstamo con información básica de la laptop.
type LaptopLoanDetail struct {
	LaptopLoan

	Brand string `db:"brand" json:"brand"`
	Model string `db:"model" json:"model"`
}
