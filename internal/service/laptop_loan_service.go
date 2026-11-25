package service

import (
	"errors"
	"time"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
)

// Errores de negocio del modulo de laptops.
var (
	ErrLaptopNotAvailable = errors.New("no hay laptops disponibles")

	ErrLaptopAlreadyBorrowed = errors.New(
		"ya tienes un prestamo activo de esta laptop",
	)
)

type LaptopLoanService struct {
	laptopLoanRepo *repository.LaptopLoanRepository
	laptopRepo     *repository.LaptopRepository
}

func NewLaptopLoanService(
	laptopLoanRepo *repository.LaptopLoanRepository,
	laptopRepo *repository.LaptopRepository,
) *LaptopLoanService {
	return &LaptopLoanService{
		laptopLoanRepo: laptopLoanRepo,
		laptopRepo:     laptopRepo,
	}
}

// Borrow registra el prestamo de una laptop.
func (s *LaptopLoanService) Borrow(
	userID int,
	req models.LaptopLoanRequest,
) (*models.LaptopLoan, error) {

	// Buscar la laptop solicitada.
	laptop, err := s.laptopRepo.FindByID(req.LaptopID)

	if err != nil {
		return nil, err
	}

	// Validar disponibilidad.
	if laptop.AvailableUnits <= 0 {
		return nil, ErrLaptopNotAvailable
	}

	// Verificar si existe un prestamo activo duplicado.
	alreadyBorrowed, err := s.laptopLoanRepo.HasActiveLoan(
		userID,
		req.LaptopID,
	)

	if err != nil {
		return nil, err
	}

	if alreadyBorrowed {
		return nil, ErrLaptopAlreadyBorrowed
	}

	// Calcular la fecha limite de devolucion.
	dueAt := time.Now().AddDate(0, 0, loanDurationDays)

	// Registrar el prestamo.
	loan, err := s.laptopLoanRepo.Create(
		userID,
		req.LaptopID,
		dueAt,
	)

	if err != nil {
		return nil, err
	}

	// Disminuir el inventario disponible.
	if err := s.laptopRepo.DecrementAvailable(req.LaptopID); err != nil {
		return nil, err
	}

	return loan, nil
}

// Return registra la devolucion de una laptop.
func (s *LaptopLoanService) Return(loanID int) error {

	// Buscar el prestamo.
	loan, err := s.laptopLoanRepo.FindByID(loanID)

	if err != nil {
		return ErrLoanNotFound
	}

	// Verificar que no haya sido devuelto.
	if loan.Status == "returned" {
		return ErrAlreadyReturned
	}

	// Registrar la devolucion.
	if err := s.laptopLoanRepo.Return(loanID); err != nil {
		return err
	}

	// Aumentar el inventario disponible.
	return s.laptopRepo.IncrementAvailable(loan.LaptopID)
}

// ListByUser consulta los prestamos de un usuario.
func (s *LaptopLoanService) ListByUser(
	userID int,
) ([]models.LaptopLoanDetail, error) {
	return s.laptopLoanRepo.FindByUser(userID)
}

// ListAll consulta todos los prestamos de laptops.
func (s *LaptopLoanService) ListAll() ([]models.LaptopLoanDetail, error) {
	return s.laptopLoanRepo.FindAll()
}
