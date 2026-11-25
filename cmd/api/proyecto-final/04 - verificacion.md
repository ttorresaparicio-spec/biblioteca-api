
# Fase 4: Verificación

## Objetivo

Verificar el correcto funcionamiento del módulo de préstamos de laptops y comprobar que la incorporación de esta nueva funcionalidad no afecte las operaciones existentes del sistema de préstamos de libros.

## Preparación del entorno

Para la verificación del sistema se utilizan las siguientes herramientas:

- Docker Desktop.
- MySQL 8.0.
- Adminer.
- Swagger UI.
- Visual Studio Code.
- API REST desarrollada en Go.

### Servicios del sistema

| Servicio | Puerto |
|---|---|
| Biblioteca API | 8082 |
| MySQL | 3308 |
| Adminer | 8083 |

La interfaz Swagger se encuentra disponible en:

http://localhost:8082/swagger/index.html

---

## Prueba 1. Contenedores Docker

**Objetivo:** verificar que los tres contenedores del proyecto estén funcionando correctamente.

**Procedimiento:** ejecutar el siguiente comando:

```powershell
docker ps
```

**Resultado esperado:** los contenedores biblioteca_final_api, biblioteca_final_mysql y biblioteca_final_adminer deben encontrarse activos.

---

## Prueba 2. Verificación de la API - Health

**Endpoint:** `GET /health`

**Objetivo:** comprobar el estado general de la aplicación.

**Resultado esperado:** respuesta HTTP 200 OK, indicando que la API se encuentra disponible.

---

## Prueba 3. Verificación de la API - Ping

**Endpoint:** `GET /api/ping`

**Objetivo:** comprobar que la API responde correctamente a las solicitudes HTTP.

**Resultado esperado:** respuesta HTTP 200 OK.

---

## Prueba 4. Registro de usuario

**Endpoint:** `POST /api/auth/register`

**Objetivo:** verificar que el sistema permite registrar nuevos usuarios.

**Resultado esperado:** respuesta HTTP 201 Created, creación del usuario con rol student y generación de un token JWT.

---

## Prueba 5. Login de usuario

**Endpoint:** `POST /api/auth/login`

**Objetivo:** comprobar el funcionamiento del mecanismo de autenticación.

**Resultado esperado:** respuesta HTTP 200 OK y generación de un token JWT válido para acceder a las rutas protegidas.

---

## Prueba 6. Consulta de libros

**Endpoint:** `GET /api/books`

**Objetivo:** verificar que el catálogo original de libros continúa funcionando después de incorporar el módulo de laptops.

**Resultado esperado:** respuesta HTTP 200 OK y listado de los libros registrados.

---

## Prueba 7. Solicitud de préstamo de libro

**Endpoint:** `POST /api/loans`

**Objetivo:** comprobar que los usuarios pueden continuar solicitando préstamos de libros.

**Resultado esperado:** respuesta HTTP 201 Created, registro del préstamo y disminución de una unidad del inventario disponible.

---

## Prueba 8. Consulta de préstamos de libros del usuario

**Endpoint:** `GET /api/loans/me`

**Objetivo:** verificar que el usuario autenticado puede consultar sus préstamos de libros.

**Resultado esperado:** respuesta HTTP 200 OK y listado de los préstamos asociados al usuario.

---

## Prueba 9. Devolución de préstamo de libro

**Endpoint:** `PUT /api/loans/{id}/return`

**Objetivo:** comprobar que el proceso original de devolución de libros continúa funcionando.

**Resultado esperado:** respuesta HTTP 200 OK, actualización del estado del préstamo a returned e incremento de la disponibilidad del libro.

---

## Prueba 10. Consulta del inventario de laptops

**Endpoint:** `GET /api/laptops`

**Objetivo:** verificar que el nuevo módulo permite consultar las laptops registradas.

**Resultado esperado:** respuesta HTTP 200 OK y listado del inventario de laptops, incluyendo marca, modelo, número de serie, unidades totales y unidades disponibles.

---

## Prueba 11. Registro de laptop

**Endpoint:** `POST /api/laptops`

**Objetivo:** comprobar que un usuario administrador puede registrar nuevas laptops en el sistema.

**Resultado esperado:** respuesta HTTP 201 Created y almacenamiento de la nueva laptop en la base de datos.

Esta operación requiere autenticación y permisos de administrador.

---

## Prueba 12. Solicitud de préstamo de laptop

**Endpoint:** `POST /api/laptop-loans`

**Objetivo:** verificar el funcionamiento del nuevo proceso de préstamo de laptops.

**Datos de ejemplo:**

```json
{
  "laptop_id": 1
}
```

**Resultado esperado:**

- Respuesta HTTP 201 Created.
- Registro del préstamo en la tabla laptop_loans.
- Estado inicial active.
- Fecha de vencimiento de 14 días.
- Disminución de una unidad del inventario disponible.

---

## Prueba 13. Consulta de préstamos de laptops del usuario

**Endpoint:** `GET /api/laptop-loans/me`

**Objetivo:** verificar que el usuario autenticado puede consultar sus préstamos de laptops.

**Resultado esperado:** respuesta HTTP 200 OK y listado de los préstamos asociados al usuario, incluyendo la información de las laptops correspondientes.

---

## Prueba 14. Devolución de laptop

**Endpoint:** `PUT /api/laptop-loans/{id}/return`

**Objetivo:** comprobar que el sistema permite realizar la devolución de una laptop prestada.

**Resultado esperado:**

- Respuesta HTTP 200 OK.
- Actualización del estado del préstamo a returned.
- Registro de la fecha de devolución.
- Incremento de una unidad del inventario disponible.

---

## Prueba 15. Validación de préstamo duplicado

**Endpoint:** `POST /api/laptop-loans`

**Objetivo:** verificar que un usuario no pueda solicitar un segundo préstamo activo de la misma laptop.

**Procedimiento:** solicitar una laptop y realizar nuevamente la misma solicitud sin devolver el préstamo anterior.

**Resultado esperado:** respuesta HTTP 400 Bad Request y mensaje de error indicando que el usuario ya tiene un préstamo activo de esa laptop.

El inventario debe permanecer sin modificaciones después de rechazar la segunda solicitud.

---

## Resumen de verificación

| N.º | Prueba | Código esperado |
|---|---|---|
| 1 | Contenedores Docker | Contenedores activos |
| 2 | Health | 200 |
| 3 | Ping | 200 |
| 4 | Registro de usuario | 201 |
| 5 | Login de usuario | 200 |
| 6 | Consulta de libros | 200 |
| 7 | Préstamo de libro | 201 |
| 8 | Consulta de préstamos de libros | 200 |
| 9 | Devolución de libro | 200 |
| 10 | Inventario de laptops | 200 |
| 11 | Registro de laptop | 201 |
| 12 | Préstamo de laptop | 201 |
| 13 | Consulta de préstamos de laptops | 200 |
| 14 | Devolución de laptop | 200 |
| 15 | Préstamo duplicado | 400 |

## Pruebas de regresión

Las pruebas relacionadas con el catálogo, préstamo, consulta y devolución de libros permiten verificar que la incorporación del módulo de laptops no afecta las funcionalidades originales de Biblioteca API.

## Consistencia de los datos

Las operaciones de préstamo y devolución de laptops utilizan transacciones SQL para mantener la consistencia entre el registro de préstamos y el inventario disponible.

En caso de producirse un error durante una operación, la transacción permite revertir los cambios realizados y evitar inconsistencias en la base de datos.

## Conclusión

El proceso de verificación permite evaluar el funcionamiento del nuevo módulo de laptops, el cumplimiento de las reglas de negocio, la autenticación de usuarios y la conservación de las funcionalidades originales.

La validación de estos resultados permite determinar si la implementación cumple con la decisión arquitectónica establecida en el ADR-001.

## Evidencias

Las capturas de pantalla correspondientes a las pruebas realizadas se encuentran en el siguiente documento:

[Ver evidencias de funcionamiento](./06%20-%20evidencias.md)

