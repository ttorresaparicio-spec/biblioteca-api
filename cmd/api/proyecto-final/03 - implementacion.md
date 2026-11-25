# Fase 3 — Implementación de la decisión

*Proyecto final · Préstamos de laptops en `biblioteca-api`*

La implementación sigue el ADR-001: incorporar un módulo independiente de laptops, conservando la arquitectura por capas y sin modificar la lógica de préstamos de libros.

> Los bloques siguientes son fragmentos transcritos de las capturas del informe (páginas 13–27), no archivos fuente completos. Las partes recortadas se documentan mediante su descripción y enlace al archivo correspondiente. Verificar los enlaces y los fragmentos contra el código definitivo al integrar los forks.

## 0. Plan de archivos nuevos y modificados

| Archivo | Tipo | Contenido |
|---|---|---|
| [migrations/002_laptops.sql](../migrations/002_laptops.sql) | Nuevo | Tablas laptops y laptop_loans |
| [internal/models/laptop.go](../internal/models/laptop.go) | Nuevo | Laptop y LaptopRequest |
| [internal/models/laptop_loan.go](../internal/models/laptop_loan.go) | Nuevo | LaptopLoan, LaptopLoanRequest y LaptopLoanDetail |
| [internal/repository/laptop_repository.go](../internal/repository/laptop_repository.go) | Nuevo | SQL del catálogo |
| [internal/repository/laptop_loan_repository.go](../internal/repository/laptop_loan_repository.go) | Nuevo | SQL de préstamos |
| [internal/service/laptop_service.go](../internal/service/laptop_service.go) | Nuevo | Reglas del catálogo |
| [internal/service/laptop_loan_service.go](../internal/service/laptop_loan_service.go) | Nuevo | Reglas de préstamos |
| [internal/handler/laptop_handler.go](../internal/handler/laptop_handler.go) | Nuevo | Endpoints del catálogo |
| [internal/handler/laptop_loan_handler.go](../internal/handler/laptop_loan_handler.go) | Nuevo | Endpoints de préstamos |
| [internal/handler/routes.go](../internal/handler/routes.go) | Editar | Registrar rutas nuevas |
| [cmd/api/main.go](../cmd/api/main.go) | Editar | Wiring de componentes |

## 1. Migración nueva — `migrations/002_laptops.sql`

Fragmento de los índices de la migración:

```sql
CREATE INDEX idx_laptop_loans_user
ON laptop_loans(user_id);

CREATE INDEX idx_laptop_loans_laptop
ON laptop_loans(laptop_id);
```

**¿Qué pasó?** Se creó un dominio independiente con dos tablas: `laptops`, que representa el catálogo de equipos, y `laptop_loans`, que representa los préstamos realizados por los usuarios. Se agregaron índices sobre `user_id` y `laptop_id`. La captura de la migración no muestra íntegramente todas las definiciones y datos iniciales.

## 2. Modelo — `internal/models/laptop.go`

```go
type Laptop struct {
    ID             int       `db:"id" json:"id"`
    Brand          string    `db:"brand" json:"brand"`
    Model          string    `db:"model" json:"model"`
    SerialNumber   string    `db:"serial_number" json:"serial_number"`
    TotalUnits     int       `db:"total_units" json:"total_units"`
    AvailableUnits int       `db:"available_units" json:"available_units"`
    CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
```

**¿Qué pasó?** Se adaptó la estructura de Book y BookRequest al dominio de laptops. LaptopRequest utiliza `binding:"required"` para validar los datos mínimos de la solicitud HTTP.

## 3. Modelo — `internal/models/laptop_loan.go`

```go
type LaptopLoan struct {
    ID         int        `db:"id" json:"id"`
    UserID     int        `db:"user_id" json:"user_id"`
    LaptopID   int        `db:"laptop_id" json:"laptop_id"`
    LoanedAt   time.Time  `db:"loaned_at" json:"loaned_at"`
    DueAt      time.Time  `db:"due_at" json:"due_at"`
    ReturnedAt *time.Time `db:"returned_at" json:"returned_at,omitempty"`
    Status     string     `db:"status" json:"status"`
}

type LaptopLoanRequest struct {
    LaptopID int `json:"laptop_id" binding:"required" example:"1"`
}

type LaptopLoanDetail struct {
    LaptopLoan
    Brand string `db:"brand" json:"brand"`
    Model string `db:"model" json:"model"`
}
```

**¿Qué pasó?** Se siguió la estructura de `models/loan.go`. LaptopLoan representa un préstamo activo o finalizado. LaptopLoanDetail incorpora información básica del equipo. La captura repite LaptopLoanRequest; aquí se presenta una sola vez.

## 4. Repository — `internal/repository/laptop_repository.go`

```go
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

func (r *LaptopRepository) IncrementAvailable(id int) error {
    query := `
        UPDATE laptops
        SET available_units = available_units + 1
        WHERE id = ?
    `
    _, err := r.db.Exec(query, id)
    return err
}
```

**¿Qué pasó?** Se creó un repositorio siguiendo el patrón de `book_repository.go`.

| Método | Responsabilidad |
|---|---|
| FindAll() | Listar el inventario |
| FindByID() | Buscar una laptop |
| Create() | Registrar una laptop |
| DecrementAvailable() | Reducir disponibilidad al prestar |
| IncrementAvailable() | Aumentar disponibilidad al devolver |

## 5. Repository — `internal/repository/laptop_loan_repository.go`

Fragmento de la consulta de préstamos por usuario:

```sql
SELECT
    ll.*,
    l.brand,
    l.model
FROM laptop_loans ll
JOIN laptops l
    ON l.id = ll.laptop_id
WHERE ll.user_id = ?
ORDER BY ll.loaned_at DESC
```

**¿Qué pasó?** Se creó el repositorio encargado de la persistencia de los préstamos.

| Método | Responsabilidad |
|---|---|
| Create() | Registrar un préstamo |
| FindByID() | Consultar un préstamo específico |
| FindByUser() | Consultar préstamos de un usuario con información de la laptop |
| FindAll() | Consultar todos los préstamos |
| MarkReturned() | Marcar devolución y registrar su fecha |
| HasActiveLoan() | Verificar si existe un préstamo activo de la misma laptop |

La captura de MarkReturned y HasActiveLoan está recortada; consultar el archivo fuente para el código completo.

## 6. Servicio — `internal/service/laptop_service.go`

```go
func (s *LaptopService) List() ([]models.Laptop, error) {
    return s.laptopRepo.FindAll()
}

func (s *LaptopService) Get(id int) (*models.Laptop, error) {
    return s.laptopRepo.FindByID(id)
}

func (s *LaptopService) Create(req models.LaptopRequest) (*models.Laptop, error) {
    laptop := &models.Laptop{
        Brand: req.Brand,
        Model: req.Model,
        SerialNumber: req.SerialNumber,
        TotalUnits: req.TotalUnits,
        AvailableUnits: req.TotalUnits,
    }
    if err := s.laptopRepo.Create(laptop); err != nil {
        return nil, err
    }
    return laptop, nil
}
```

**¿Qué pasó?** Se creó el servicio del catálogo siguiendo la estructura de BookService. List consulta el inventario, Get obtiene un equipo por ID y Create registra equipos nuevos. Se conserva la separación entre LaptopService y LaptopLoanService.

## 7. Servicio — `internal/service/laptop_loan_service.go`

Fragmentos de Borrow:

```go
if laptop.AvailableUnits <= 0 {
    return nil, ErrLaptopNotAvailable
}
```

```go
dueAt := time.Now().AddDate(0, 0, loanDurationDays)
```

**¿Qué pasó?** Borrow busca la laptop, valida disponibilidad, verifica préstamos duplicados, calcula el vencimiento, registra el préstamo y actualiza el inventario. Return busca el préstamo, verifica que no esté devuelto, lo marca como retornado e incrementa la disponibilidad. ListByUser consulta los préstamos del usuario y ListAll permite la consulta administrativa.

**Decisión arquitectónica:** se reutiliza `loanDurationDays` (14 días), sin modificar LoanService, LoanRepository ni BookRepository.

## 8. Handler — `internal/handler/laptop_handler.go`

Fragmento de validación de la solicitud de creación:

```go
var req models.LaptopRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

**¿Qué pasó?** Se creó el controlador del catálogo. List utiliza LaptopService.List y Create utiliza LaptopService.Create. Se conserva el patrón de validación mediante ShouldBindJSON.

## 9. Handler — `internal/handler/laptop_loan_handler.go`

Fragmento de validación de la solicitud de préstamo:

```go
var req models.LaptopLoanRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

**¿Qué pasó?** Se creó el controlador de préstamos de laptops con los siguientes endpoints:

| Endpoint | Funcionalidad |
|---|---|
| POST /api/laptop-loans | Solicitar laptop |
| GET /api/laptop-loans/me | Ver mis préstamos |
| PUT /api/laptop-loans/:id/return | Devolver laptop |
| GET /api/laptop-loans | Ver todos los préstamos (admin) |

## 10. Rutas — `internal/handler/routes.go` (edición)

### Rutas protegidas

```go
// Catálogo de laptops
protected.GET("/laptops", laptopHandler.List)

// Préstamos de laptops
protected.POST("/laptop-loans", laptopLoanHandler.Borrow)
protected.GET("/laptop-loans/me", laptopLoanHandler.MyLoans)
protected.PUT(
    "/laptop-loans/:id/return",
    laptopLoanHandler.Return,
)
```

### Rutas de administrador

```go
admin.POST("/laptops", laptopHandler.Create)
admin.GET("/laptop-loans", laptopLoanHandler.AllLoans)
```

**¿Qué pasó?** Se registraron las rutas del módulo reutilizando JWT, AuthRequired y AdminRequired. Las rutas protegidas requieren autenticación y las administrativas requieren permisos de administrador. No fue necesario crear middleware nuevo.

## 11. Hook de integración

**No requerido.** A diferencia del ejemplo de reservas, el módulo de laptops no requiere comunicación directa con los servicios de libros o préstamos. Las reglas del nuevo dominio permanecen encapsuladas en LaptopService y LaptopLoanService. No fue necesario modificar LoanService, BookService, LoanRepository ni BookRepository.

## 12. Wiring — `cmd/api/main.go` (edición)

### Repositorios

```go
laptopRepo := repository.NewLaptopRepository(database)
laptopLoanRepo := repository.NewLaptopLoanRepository(database)
```

### Servicios

```go
laptopService := service.NewLaptopService(laptopRepo)
laptopLoanService := service.NewLaptopLoanService(
    laptopLoanRepo,
    laptopRepo,
)
```

### Handlers

```go
laptopHandler := handler.NewLaptopHandler(laptopService)
laptopLoanHandler := handler.NewLaptopLoanHandler(laptopLoanService)
```

### Registro de rutas

```go
handler.SetupRoutes(
    router,
    authHandler,
    bookHandler,
    loanHandler,
    laptopHandler,
    laptopLoanHandler,
    cfg.JWTSecret,
)
```

**¿Qué pasó?** Se conserva el patrón de composición del proyecto: Repository → Service → Handler → Routes. Los nuevos componentes se inicializan en main.go y se registran en SetupRoutes.

## 13. Cierre

El informe registra que se ejecutó `docker compose up --build` sin errores. Esta afirmación corresponde al informe y no constituye una nueva comprobación de ejecución. Las evidencias funcionales y de regresión se presentan en la [Fase 4 — Verificación](verificacion.md).


