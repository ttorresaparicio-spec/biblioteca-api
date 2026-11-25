package service

import (
	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
)

type LaptopService struct {
	laptopRepo *repository.LaptopRepository
}

func NewLaptopService(laptopRepo *repository.LaptopRepository) *LaptopService {
	return &LaptopService{laptopRepo: laptopRepo}
}

// List devuelve todas las laptops registradas.
func (s *LaptopService) List() ([]models.Laptop, error) {
	return s.laptopRepo.FindAll()
}

// Get obtiene una laptop por su identificador.
func (s *LaptopService) Get(id int) (*models.Laptop, error) {
	return s.laptopRepo.FindByID(id)
}

// Create registra una nueva laptop.
func (s *LaptopService) Create(req models.LaptopRequest) (*models.Laptop, error) {
	laptop := &models.Laptop{
		Brand:          req.Brand,
		Model:          req.Model,
		SerialNumber:   req.SerialNumber,
		TotalUnits:     req.TotalUnits,
		AvailableUnits: req.TotalUnits,
	}

	if err := s.laptopRepo.Create(laptop); err != nil {
		return nil, err
	}

	return laptop, nil
}
