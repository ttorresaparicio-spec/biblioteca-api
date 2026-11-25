# Fase 4 — Verificación

*Proyecto final · Préstamos de laptops en biblioteca-api*

## 1. Objetivo

Verificar el funcionamiento de la nueva funcionalidad de préstamos de laptops y comprobar que su incorporación no afecte el módulo existente de préstamos de libros.

## 2. Entorno de pruebas

- Go
- Gin
- MySQL 8
- Docker
- Docker Compose
- Adminer
- Cliente HTTP

## 3. Compilación y ejecución

Comandos de verificación:

```bash
docker compose up -d --build
docker compose ps
```

**Evidencia:** Pendiente de incorporar.

## 4. Checklist de pruebas funcionales

| ID | Prueba | Resultado esperado | Estado |
|---|---|---|---|
| P01 | Consultar catálogo de laptops | Consulta exitosa | Pendiente |
| P02 | Registrar una laptop | Registro exitoso | Pendiente |
| P03 | Solicitar préstamo disponible | Préstamo registrado | Pendiente |
| P04 | Solicitar préstamo duplicado | Solicitud rechazada | Pendiente |
| P05 | Solicitar laptop sin disponibilidad | Solicitud rechazada | Pendiente |
| P06 | Consultar préstamos personales | Consulta exitosa | Pendiente |
| P07 | Devolver una laptop | Devolución registrada | Pendiente |
| P08 | Repetir una devolución | Solicitud rechazada | Pendiente |
| P09 | Consultar todos los préstamos | Consulta exitosa | Pendiente |
| P10 | Acceder sin autenticación | Acceso rechazado | Pendiente |
| P11 | Acceder sin permisos administrativos | Acceso rechazado | Pendiente |

## 5. Verificación del inventario

- [ ] Disminución de unidades disponibles al prestar.
- [ ] Incremento de unidades disponibles al devolver.
- [ ] Registro de la fecha de vencimiento.
- [ ] Registro de la fecha de devolución.
- [ ] Actualización del estado del préstamo.

## 6. Pruebas de regresión

| ID | Prueba | Estado |
|---|---|---|
| R01 | Consultar catálogo de libros | Pendiente |
| R02 | Solicitar préstamo de libro | Pendiente |
| R03 | Consultar préstamos de libros | Pendiente |
| R04 | Devolver un libro | Pendiente |

## 7. Evidencias

Como respaldo del proceso de verificación, se adjuntan las evidencias correspondientes a las pruebas funcionales, las validaciones de las reglas de negocio y las pruebas de regresión realizadas.

Las capturas de pantalla y los resultados obtenidos se encuentran organizados en el apartado de evidencias, ubicado al final de la documentación del proyecto.

[Consultar evidencias de verificación](06-evidencias.md)


## 8. Conclusión

Pendiente de completar después de verificar los resultados de las pruebas funcionales y de regresión.
