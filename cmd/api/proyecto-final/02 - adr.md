# Fase 2 — Decisión arquitectónica

## ADR-001 — Incorporación de un módulo de préstamos de laptops

**Estado:** Aceptado

**Fecha:** Noviembre de 2026

## 1. Contexto

La aplicación biblioteca-api permite gestionar préstamos de libros mediante una arquitectura por capas.

El análisis arquitectónico identificó dependencias directas entre el módulo de préstamos y la entidad Book, tanto en los modelos como en los repositorios, servicios y estructura de la base de datos.

Se requiere incorporar una funcionalidad que permita gestionar préstamos de laptops sin afectar el funcionamiento existente.

## 2. Alternativas consideradas

- Implementar un módulo independiente de laptops.
- Generalizar el módulo de préstamos mediante una entidad de recursos.
- Utilizar una tabla única de préstamos con un tipo de recurso.

Las alternativas se encuentran detalladas en el [Fase 1 — Análisis arquitectónico](01%20-%20analisis.md)

## 3. Decisión

Se decide implementar un módulo independiente para la gestión y préstamo de laptops, compuesto por:

- Tabla laptops.
- Tabla laptop_loans.
- Modelos específicos.
- Repositorios específicos.
- Servicios específicos.
- Handlers específicos.
- Nuevos endpoints.

La nueva funcionalidad conserva la arquitectura por capas y reutiliza los mecanismos existentes de autenticación y autorización.

## 4. Justificación

La alternativa seleccionada permite reducir el impacto sobre los componentes existentes, facilitar las pruebas y minimizar el riesgo de regresiones.

No se requiere modificar la lógica de negocio del módulo de préstamos de libros ni generalizar las entidades existentes.

## 5. Consecuencias positivas

- Menor riesgo de regresiones.
- Conservación de la arquitectura existente.
- Separación de responsabilidades.
- Independencia entre los módulos.
- Facilidad de mantenimiento y pruebas.

## 6. Consecuencias negativas

- Duplicación parcial de lógica de negocio.
- Menor reutilización de código.
- Posible incremento de componentes si se incorporan nuevos tipos de recursos.

## 7. Trade-off principal

Se prioriza la estabilidad del sistema y el bajo riesgo de implementación sobre la reutilización máxima de código y la extensibilidad futura.

## 8. Hook de integración

No requerido.

El módulo de préstamos de laptops funciona de manera independiente al módulo de préstamos de libros, compartiendo únicamente la infraestructura y los mecanismos de autenticación y autorización existentes.

