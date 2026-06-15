# Taller 3 — Algoritmos y Estructuras de Datos
## Estructuras lineales y no lineales aplicadas a datasets reales
**Lenguaje:** Go (Golang) 1.21+  
**Curso:** Algoritmos y Estructuras de Datos

---

## 📁 Estructura del repositorio

```
taller3/
├── go.mod
├── README.md
├── pilas/              ← Ejercicio 4.1: Stack + Stock Span
│   ├── pilas.go
│   ├── pilas_test.go
│   └── main/
│       └── main.go
└── diagramas/          ← Diagramas de funciones (PNG/PDF)
```

---

## ▶️ Instrucciones de ejecución

### Requisitos previos
- Go 1.21 o superior instalado → [go.dev/dl](https://go.dev/dl)
- Clonar el repositorio e ingresar a la carpeta raíz

```bash
git clone <url-del-repositorio>
cd taller3
```

### Ejercicio 1 — Pilas: Stock Span

**Correr el programa principal con el dataset:**
```bash
go run ./pilas/main <ruta_al_archivo.txt> [N]
```

Ejemplo con los primeros y últimos 10 días:
```bash
go run ./pilas/main adra.us.txt 10
```

**Correr los tests unitarios:**
```bash
go test ./pilas/... -v
```

**Correr los tests con cobertura:**
```bash
go test ./pilas/... -v -cover
```

**Correr los benchmarks:**
```bash
go test ./pilas/... -bench="." -benchmem -benchtime=3s -run=^$

go test ./pilas/... -bench=. -benchmem -benchtime=3s -run=^$


```

---

## 🎥 Video explicativo

> 📺 **Enlace al video en YouTube:** [INSERTAR ENLACE AQUÍ]

El video incluye:
- Presentación de los integrantes del grupo
- Explicación de la estructura Stack genérica
- Demostración del algoritmo Stock Span en ejecución
- Análisis de resultados de performance

---

## 📊 Origen del dataset

### Ejercicio 1 — Pilas (Stock Span)
**Dataset:** Huge Stock Market Dataset  
**Autor:** Boris Marjanovic  
**Fuente:** [Kaggle — Price Volume Data for All US Stocks & ETFs](https://www.kaggle.com/datasets/borismarjanovic/price-volume-data-for-all-us-stocks-etfs)  
**Archivo utilizado:** `adra.us.txt` (columna `Close`)  
**Total de registros procesados:** 2,177 días de cotización  

> ⚠️ El dataset no se sube al repositorio por su tamaño.  
> Descargarlo desde el enlace de Kaggle y colocarlo en la raíz del proyecto.

---

## 👥 Integrantes


---
