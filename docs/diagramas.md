# Fase 5 — Diagramas

*Proyecto final · Feature: **Préstamos de laptops** en `biblioteca-api`*

> Este documento presenta los diagramas correspondientes a la Fase 5 del proyecto, elaborados en Mermaid y coherentes con la decisión arquitectónica ADR-001: implementación de un módulo independiente para la gestión y préstamo de laptops, sin afectar el módulo existente de préstamos de libros.

---

## 1. C4 Nivel 1 — Contexto

**Pregunta:** ¿Quiénes interactúan con el sistema y qué esperan de él?

```mermaid
flowchart LR
    E["Estudiante (Persona)"]
    A["Bibliotecario / Administrador (Persona)"]
    C["Cliente de la biblioteca"]
    API["biblioteca-api (Sistema)"]
    BD[("MySQL 8")]

    E -->|Consulta, solicita y devuelve laptops| C
    A -->|Gestiona catalogo y prestamos| C
    C -->|REST + JWT| API
    API -->|SQL| BD
```

**¿Qué muestra?** Representa la interacción de los estudiantes y administradores con el sistema. Los estudiantes pueden consultar el catálogo, solicitar préstamos y devolver laptops, mientras que los administradores pueden registrar equipos y consultar los préstamos realizados. La comunicación se realiza mediante solicitudes REST protegidas por JWT y los datos se almacenan en MySQL.

---

## 2. C4 Nivel 2 — Contenedores

**Pregunta:** ¿De qué bloques ejecutables está compuesto el sistema?

```mermaid
flowchart LR
    C["Cliente HTTP"]
    A["Administrador de BD"]

    subgraph Docker["Docker Compose"]
        API["biblioteca-api - Go + Gin - Puerto 8080"]
        BD[("MySQL 8 - Puerto interno 3306")]
        ADM["Adminer - Puerto externo 8081"]

        API -->|TCP / SQL| BD
        ADM -->|SQL| BD
    end

    C -->|HTTP / JSON :8080| API
    A -->|HTTP :8081| ADM
```

**¿Qué muestra?** Los contenedores que conforman la aplicación: la API desarrollada en Go con Gin, MySQL como sistema gestor de base de datos y Adminer como herramienta de administración. La nueva funcionalidad de préstamos de laptops no requiere contenedores adicionales, ya que se incorpora dentro de la API y utiliza la base de datos existente.

MySQL utiliza el puerto interno 3306 y el puerto externo 3307, configurado para evitar conflictos con el servicio local de MySQL.

---

## 3. Diagrama de componentes

**Pregunta:** ¿Qué módulos internos tiene la API y dónde se encuentra la nueva funcionalidad?

```mermaid
flowchart TD
    R["Routes - JWT y permisos"]

    subgraph Existente["Modulo existente - Libros"]
        BH[BookHandler]
        LH[LoanHandler]
        BS[BookService]
        LS[LoanService]
        BR[BookRepository]
        LR[LoanRepository]

        BH --> BS --> BR
        LH --> LS --> LR
        LS --> BR
    end

    subgraph Nuevo["Nuevo modulo - Laptops"]
        LPH[LaptopHandler]
        LLH[LaptopLoanHandler]
        LPS[LaptopService]
        LLS[LaptopLoanService]
        LPR[LaptopRepository]
        LLR[LaptopLoanRepository]

        LPH --> LPS --> LPR
        LLH --> LLS --> LLR
        LLS --> LPR
    end

    R --> BH
    R --> LH
    R --> LPH
    R --> LLH

    BR --> BD[(MySQL)]
    LR --> BD
    LPR --> BD
    LLR --> BD

    style Nuevo fill:#163c2b,stroke:#45a875,color:#ffffff
```

**¿Qué muestra?** La organización interna de la API mediante las capas Handler, Service y Repository. El módulo de laptops funciona de manera paralela e independiente al módulo de libros.

Los nuevos componentes reutilizan los mecanismos de autenticación, autorización y conexión a la base de datos existentes, sin establecer dependencias directas con el módulo de préstamos de libros.

Esta estructura respeta la decisión arquitectónica ADR-001 y reduce el riesgo de regresiones.

---

## 4. Diagrama de secuencia — Solicitar préstamo de una laptop

**Pregunta:** ¿Qué ocurre, en orden, cuando un estudiante solicita el préstamo de una laptop?

```mermaid
sequenceDiagram
    actor U as Estudiante
    participant R as Routes / JWT
    participant H as LaptopLoanHandler
    participant S as LaptopLoanService
    participant LP as LaptopRepository
    participant LL as LaptopLoanRepository
    participant BD as MySQL

    U->>R: POST /api/laptop-loans
    R->>R: Validar JWT
    R->>H: Borrow()
    H->>H: Validar JSON
    H->>S: Borrow(userID, req)
    S->>LP: FindByID(laptopID)
    LP->>BD: Consultar laptop
    BD-->>LP: Resultado
    LP-->>S: Laptop o error

    alt Laptop no encontrada
        S-->>H: Error
        H-->>U: HTTP 400
    else Laptop encontrada
        alt Sin unidades disponibles
            S-->>H: ErrLaptopNotAvailable
            H-->>U: HTTP 400
        else Unidades disponibles
            S->>LL: HasActiveLoan(userID, laptopID)
            LL->>BD: Consultar prestamos activos
            BD-->>LL: Resultado
            LL-->>S: Existe prestamo activo

            alt Prestamo duplicado
                S-->>H: ErrLaptopAlreadyBorrowed
                H-->>U: HTTP 400
            else Sin prestamo duplicado
                S->>S: Calcular vencimiento de 14 dias
                S->>LL: Create(userID, laptopID, dueAt)
                LL->>BD: INSERT INTO laptop_loans
                BD-->>LL: Prestamo registrado
                LL-->>S: LaptopLoan
                S->>LP: DecrementAvailable(laptopID)
                LP->>BD: Actualizar available_units
                BD-->>LP: Inventario actualizado
                LP-->>S: Operacion exitosa
                S-->>H: LaptopLoan
                H-->>U: HTTP 201 Created
            end
        end
    end
```

**¿Qué muestra?** El flujo de comunicación entre los componentes del sistema durante la solicitud de un préstamo de laptop.

El servicio valida la existencia del equipo, verifica que existan unidades disponibles y comprueba que el usuario no tenga un préstamo activo de la misma laptop.

Si las validaciones son exitosas, se registra el préstamo con una duración de 14 días y se disminuye la cantidad de unidades disponibles en el inventario.

Las estructuras alt representan las condiciones y validaciones de negocio implementadas en LaptopLoanService.

---

## 5. Diagrama de secuencia — Devolución de una laptop

**Pregunta:** ¿Cómo se registra la devolución de una laptop y se actualiza su disponibilidad?

```mermaid
sequenceDiagram
    actor U as Usuario
    participant R as Routes / JWT
    participant H as LaptopLoanHandler
    participant S as LaptopLoanService
    participant LL as LaptopLoanRepository
    participant LP as LaptopRepository
    participant BD as MySQL

    U->>R: PUT /api/laptop-loans/:id/return
    R->>R: Validar JWT
    R->>H: Return()
    H->>H: Validar ID
    H->>S: Return(loanID)
    S->>LL: FindByID(loanID)
    LL->>BD: Consultar laptop_loans
    BD-->>LL: Resultado
    LL-->>S: Prestamo o error

    alt Prestamo no encontrado
        S-->>H: ErrLoanNotFound
        H-->>U: HTTP 400
    else Prestamo encontrado
        alt Prestamo ya devuelto
            S-->>H: ErrAlreadyReturned
            H-->>U: HTTP 400
        else Prestamo activo
            S->>LL: MarkReturned(loanID)
            LL->>BD: Actualizar estado y fecha
            BD-->>LL: Prestamo actualizado
            LL-->>S: Operacion exitosa

            S->>LP: IncrementAvailable(laptopID)
            LP->>BD: Incrementar available_units
            BD-->>LP: Inventario actualizado
            LP-->>S: Operacion exitosa

            S-->>H: Devolucion exitosa
            H-->>U: HTTP 200 OK
        end
    end
```

**¿Qué muestra?** El proceso de devolución de una laptop, incluyendo la validación de la existencia del préstamo y la comprobación de que no haya sido devuelto anteriormente.

Una vez superadas las validaciones, el sistema actualiza el estado del préstamo, registra la fecha de devolución e incrementa la disponibilidad de la laptop.

Este proceso se realiza exclusivamente dentro del módulo de laptops, sin modificar ni utilizar los servicios de préstamos de libros.

**Hook de integración:** No requerido. La independencia entre ambos módulos permite conservar el funcionamiento del sistema existente, respetando la decisión arquitectónica ADR-001.

---


