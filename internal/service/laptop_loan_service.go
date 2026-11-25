package service

import (
	"time"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
)

type LaptopLoanService struct {
	laptopLoanRepo *repository.LaptopLoanRepository
	laptopRepo     *repository.LaptopRepository
}

// Constructor del servicio de prestamos de laptops.
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

	// Calcular la fecha limite de devolucion.
	dueAt := time.Now().AddDate(0, 0, loanDurationDays)

	// Registrar el prestamo mediante una transaccion.
	return s.laptopLoanRepo.BorrowAtomic(
		userID,
		req.LaptopID,
		dueAt,
	)
}

// Return registra la devolucion de una laptop.
func (s *LaptopLoanService) Return(loanID int) error {

	return s.laptopLoanRepo.ReturnAtomic(loanID)
}

// ListByUser consulta los prestamos de un usuario.
func (s *LaptopLoanService) ListByUser(
	userID int,
) ([]models.LaptopLoanDetail, error) {

	return s.laptopLoanRepo.FindByUser(userID)
}

// ListAll consulta todos los prestamos de laptops.
func (s *LaptopLoanService) ListAll() (
	[]models.LaptopLoanDetail, error,
) {

	return s.laptopLoanRepo.FindAll()
}
