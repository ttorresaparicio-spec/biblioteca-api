package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/service"
)

type LaptopHandler struct {
	laptopService *service.LaptopService
}

func NewLaptopHandler(laptopService *service.LaptopService) *LaptopHandler {
	return &LaptopHandler{
		laptopService: laptopService,
	}
}

// List godoc
// @Summary Listar laptops
// @Tags laptops
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Laptop
// @Router /laptops [get]
func (h *LaptopHandler) List(c *gin.Context) {
	laptops, err := h.laptopService.List()

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "no se pudieron obtener las laptops"},
		)
		return
	}

	c.JSON(http.StatusOK, laptops)
}

// Create godoc
// @Summary Crear laptop
// @Tags laptops
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.LaptopRequest true "Datos de la laptop"
// @Success 201 {object} models.Laptop
// @Router /laptops [post]
func (h *LaptopHandler) Create(c *gin.Context) {
	var req models.LaptopRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	laptop, err := h.laptopService.Create(req)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusCreated, laptop)
}
