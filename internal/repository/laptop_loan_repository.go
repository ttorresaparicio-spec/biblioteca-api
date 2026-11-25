package repository

import (
	"errors"
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

// BorrowAtomic registra un prestamo y actualiza el inventario
// dentro de una misma transaccion.
func (r *LaptopLoanRepository) BorrowAtomic(
	userID, laptopID int,
	dueAt time.Time,
) (*models.LaptopLoan, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Bloquear la laptop mientras se realiza el prestamo.
	var laptop models.Laptop

	err = tx.Get(&laptop,
		`SELECT * FROM laptops WHERE id = ? FOR UPDATE`,
		laptopID,
	)
	if err != nil {
		return nil, err
	}

	// Verificar disponibilidad.
	if laptop.AvailableUnits <= 0 {
		return nil, errors.New("no hay laptops disponibles")
	}

	// Verificar prestamos activos duplicados.
	var count int

	err = tx.Get(&count, `
        SELECT COUNT(*)
        FROM laptop_loans
        WHERE user_id = ?
        AND laptop_id = ?
        AND status = 'active'
    `, userID, laptopID)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New(
			"ya tienes un prestamo activo de esta laptop",
		)
	}

	// Registrar el prestamo.
	result, err := tx.Exec(`
        INSERT INTO laptop_loans
        (user_id, laptop_id, due_at, status)
        VALUES (?, ?, ?, 'active')
    `, userID, laptopID, dueAt)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Disminuir el inventario.
	result, err = tx.Exec(`
        UPDATE laptops
        SET available_units = available_units - 1
        WHERE id = ?
        AND available_units > 0
    `, laptopID)

	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected != 1 {
		return nil, errors.New("no se pudo actualizar el inventario")
	}

	// Consultar el prestamo creado.
	var loan models.LaptopLoan

	err = tx.Get(&loan,
		`SELECT * FROM laptop_loans WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}

	// Confirmar la transaccion.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &loan, nil
}

// ReturnAtomic registra la devolucion y recupera
// una unidad del inventario mediante una transaccion.
func (r *LaptopLoanRepository) ReturnAtomic(loanID int) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Consultar el prestamo.
	var loan models.LaptopLoan

	err = tx.Get(&loan,
		`SELECT * FROM laptop_loans WHERE id = ?`,
		loanID,
	)
	if err != nil {
		return err
	}

	// Bloquear la laptop para controlar operaciones simultaneas.
	var laptop models.Laptop

	err = tx.Get(&laptop,
		`SELECT * FROM laptops WHERE id = ? FOR UPDATE`,
		loan.LaptopID,
	)
	if err != nil {
		return err
	}

	// Actualizar solamente prestamos activos.
	result, err := tx.Exec(`
        UPDATE laptop_loans
        SET status = 'returned',
            returned_at = NOW()
        WHERE id = ?
        AND status = 'active'
    `, loanID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("el prestamo ya fue devuelto")
	}

	// Recuperar una unidad sin superar el inventario total.
	result, err = tx.Exec(`
        UPDATE laptops
        SET available_units = available_units + 1
        WHERE id = ?
        AND available_units < total_units
    `, loan.LaptopID)

	if err != nil {
		return err
	}

	affected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("no se pudo recuperar el inventario")
	}

	// Confirmar la transaccion.
	return tx.Commit()
}
