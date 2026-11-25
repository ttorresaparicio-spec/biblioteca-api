package repository

import (
	"time"

	"github.com/jmoiron/sqlx"

	"biblioteca-api/internal/models"
)

type LaptopLoanRepository struct {
	db *sqlx.DB
}

func NewLaptopLoanRepository(db *sqlx.DB) *LaptopLoanRepository {
	return &LaptopLoanRepository{db: db}
}

// Create registra un nuevo préstamo de laptop.
func (r *LaptopLoanRepository) Create(userID, laptopID int, dueAt time.Time) (*models.LaptopLoan, error) {
	query := `
		INSERT INTO laptop_loans (
			user_id,
			laptop_id,
			due_at,
			status
		)
		VALUES (?, ?, ?, 'active')
	`

	result, err := r.db.Exec(query, userID, laptopID, dueAt)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	return r.FindByID(int(id))
}

// FindByID busca un préstamo específico.
func (r *LaptopLoanRepository) FindByID(id int) (*models.LaptopLoan, error) {
	var loan models.LaptopLoan

	query := `SELECT * FROM laptop_loans WHERE id = ?`

	if err := r.db.Get(&loan, query, id); err != nil {
		return nil, err
	}

	return &loan, nil
}

// FindByUser devuelve los préstamos de laptops de un usuario.
func (r *LaptopLoanRepository) FindByUser(userID int) ([]models.LaptopLoanDetail, error) {
	query := `
		SELECT
			ll.*,
			l.brand,
			l.model
		FROM laptop_loans ll
		JOIN laptops l
			ON l.id = ll.laptop_id
		WHERE ll.user_id = ?
		ORDER BY ll.loaned_at DESC
	`

	var loans []models.LaptopLoanDetail

	if err := r.db.Select(&loans, query, userID); err != nil {
		return nil, err
	}

	return loans, nil
}

// FindAll devuelve todos los préstamos de laptops.
func (r *LaptopLoanRepository) FindAll() ([]models.LaptopLoanDetail, error) {
	query := `
		SELECT
			ll.*,
			l.brand,
			l.model
		FROM laptop_loans ll
		JOIN laptops l
			ON l.id = ll.laptop_id
		ORDER BY ll.loaned_at DESC
	`

	var loans []models.LaptopLoanDetail

	if err := r.db.Select(&loans, query); err != nil {
		return nil, err
	}

	return loans, nil
}

// Return registra la devolución de una laptop.
func (r *LaptopLoanRepository) Return(id int) error {
	query := `
		UPDATE laptop_loans
		SET status = 'returned',
			returned_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}

// HasActiveLoan verifica si el usuario ya tiene
// un préstamo activo de la misma laptop.
func (r *LaptopLoanRepository) HasActiveLoan(userID, laptopID int) (bool, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM laptop_loans
		WHERE user_id = ?
		AND laptop_id = ?
		AND status = 'active'
	`

	if err := r.db.Get(&count, query, userID, laptopID); err != nil {
		return false, err
	}

	return count > 0, nil
}
