package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/service"
)

type LaptopLoanHandler struct {
	laptopLoanService *service.LaptopLoanService
}

func NewLaptopLoanHandler(
	laptopLoanService *service.LaptopLoanService,
) *LaptopLoanHandler {
	return &LaptopLoanHandler{
		laptopLoanService: laptopLoanService,
	}
}

// Borrow godoc
// @Summary Solicitar prestamo de laptop
// @Tags laptop-loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.LaptopLoanRequest true "ID de la laptop"
// @Success 201 {object} models.LaptopLoan
// @Router /laptop-loans [post]
func (h *LaptopLoanHandler) Borrow(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.LaptopLoanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	loan, err := h.laptopLoanService.Borrow(userID, req)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusCreated, loan)
}

// MyLoans godoc
// @Summary Mis prestamos de laptop
// @Tags laptop-loans
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.LaptopLoanDetail
// @Router /laptop-loans/me [get]
func (h *LaptopLoanHandler) MyLoans(c *gin.Context) {
	userID := c.GetInt("user_id")

	loans, err := h.laptopLoanService.ListByUser(userID)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "no se pudieron obtener los prestamos"},
		)
		return
	}

	c.JSON(http.StatusOK, loans)
}

// Return godoc
// @Summary Devolver laptop
// @Tags laptop-loans
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del prestamo"
// @Success 200 {object} map[string]string
// @Router /laptop-loans/{id}/return [put]
func (h *LaptopLoanHandler) Return(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "ID invalido"},
		)
		return
	}

	// Comprobar que el prestamo pertenece al usuario.
	userID := c.GetInt("user_id")

	loans, err := h.laptopLoanService.ListByUser(userID)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "no se pudieron consultar los prestamos"},
		)
		return
	}

	authorized := false

	for _, loan := range loans {
		if loan.ID == id {
			authorized = true
			break
		}
	}

	if !authorized {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "no tienes permiso para devolver este prestamo"},
		)
		return
	}

	if err := h.laptopLoanService.Return(id); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"message": "laptop devuelta correctamente"},
	)
}

// AllLoans godoc
// @Summary Ver todos los prestamos de laptops
// @Tags laptop-loans
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.LaptopLoanDetail
// @Router /laptop-loans [get]
func (h *LaptopLoanHandler) AllLoans(c *gin.Context) {
	loans, err := h.laptopLoanService.ListAll()

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "no se pudieron obtener los prestamos"},
		)
		return
	}

	c.JSON(http.StatusOK, loans)
}
