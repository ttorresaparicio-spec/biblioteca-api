# Diagramas de arquitectura

## 1. C4 Nivel 1 - Contexto

```mermaid
graph LR
    A[Estudiante] --> C[App Web]
    B[Bibliotecario] --> C
    C --> D[biblioteca-api]
    D --> E[(MySQL 8)]
```
