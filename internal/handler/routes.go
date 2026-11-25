package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"biblioteca-api/internal/middleware"
)

// SetupRoutes registra todas las rutas de la API sobre el router de Gin.
func SetupRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	bookHandler *BookHandler,
	loanHandler *LoanHandler,
	laptopHandler *LaptopHandler,
	laptopLoanHandler *LaptopLoanHandler,
	jwtSecret string,
) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := router.Group("/api")

	// Rutas públicas de autenticación.
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Rutas protegidas mediante JWT.
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtSecret))
	{
		protected.GET("/auth/me", authHandler.Me)

		// Catálogo de libros.
		protected.GET("/books", bookHandler.List)
		protected.GET("/books/:id", bookHandler.Get)

		// Préstamos de libros.
		protected.POST("/loans", loanHandler.Borrow)
		protected.GET("/loans/me", loanHandler.MyLoans)
		protected.PUT("/loans/:id/return", loanHandler.Return)

		// Catálogo de laptops.
		protected.GET("/laptops", laptopHandler.List)

		// Préstamos de laptops.
		protected.POST("/laptop-loans", laptopLoanHandler.Borrow)
		protected.GET("/laptop-loans/me", laptopLoanHandler.MyLoans)
		protected.PUT("/laptop-loans/:id/return", laptopLoanHandler.Return)
	}

	// Rutas exclusivas para administradores.
	admin := api.Group("")
	admin.Use(
		middleware.AuthRequired(jwtSecret),
		middleware.AdminRequired(),
	)
	{
		// Administración de libros.
		admin.POST("/books", bookHandler.Create)
		admin.PUT("/books/:id", bookHandler.Update)
		admin.DELETE("/books/:id", bookHandler.Delete)
		admin.GET("/loans", loanHandler.AllLoans)

		// Administración de laptops.
		admin.POST("/laptops", laptopHandler.Create)
		admin.GET("/laptop-loans", laptopLoanHandler.AllLoans)
	}
}
