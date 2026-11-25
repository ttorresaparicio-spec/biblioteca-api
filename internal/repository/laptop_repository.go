package repository

import (
	"github.com/jmoiron/sqlx"

	"biblioteca-api/internal/models"
)

type LaptopRepository struct {
	db *sqlx.DB
}

func NewLaptopRepository(db *sqlx.DB) *LaptopRepository {
	return &LaptopRepository{db: db}
}

// FindAll devuelve todas las laptops registradas.
func (r *LaptopRepository) FindAll() ([]models.Laptop, error) {
	query := `SELECT * FROM laptops ORDER BY brand ASC, model ASC`

	var laptops []models.Laptop

	if err := r.db.Select(&laptops, query); err != nil {
		return nil, err
	}

	return laptops, nil
}

// FindByID busca una laptop por su identificador.
func (r *LaptopRepository) FindByID(id int) (*models.Laptop, error) {
	var laptop models.Laptop

	query := `SELECT * FROM laptops WHERE id = ?`

	if err := r.db.Get(&laptop, query, id); err != nil {
		return nil, err
	}

	return &laptop, nil
}

// Create registra una nueva laptop.
func (r *LaptopRepository) Create(l *models.Laptop) error {
	query := `
		INSERT INTO laptops (
			brand,
			model,
			serial_number,
			total_units,
			available_units
		)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		l.Brand,
		l.Model,
		l.SerialNumber,
		l.TotalUnits,
		l.AvailableUnits,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	l.ID = int(id)

	return nil
}

// DecrementAvailable reduce una unidad disponible al prestar una laptop.
func (r *LaptopRepository) DecrementAvailable(id int) error {
	query := `
		UPDATE laptops
		SET available_units = available_units - 1
		WHERE id = ?
		AND available_units > 0
	`

	_, err := r.db.Exec(query, id)

	return err
}

// IncrementAvailable aumenta una unidad disponible al devolver una laptop.
func (r *LaptopRepository) IncrementAvailable(id int) error {
	query := `
		UPDATE laptops
		SET available_units = available_units + 1
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
