# Taller 3 — Ejercicio 4.3: Caché LRU con Listas Enlazadas

## Descripción

Implementación de una caché **LRU (Least Recently Used)** en Go, usando una
lista doblemente enlazada (con nodos ficticios `head`/`tail`) combinada con
un mapa para acceso en O(1).

La caché se evalúa con una secuencia real de accesos extraída del dataset
**MovieLens**: se ordena `ratings.csv` por `timestamp` y se toma la columna
`movieId` como la secuencia de "películas pedidas" en orden cronológico.
Se calcula el **hit ratio** (aciertos / accesos totales) para distintos
tamaños de caché.

## Dataset

- **Fuente**: MovieLens (GroupLens) — `ratings.csv`
- **Link**: https://grouplens.org/datasets/movielens/
- **Columnas usadas**: `movieId`, `timestamp`
- El archivo `ratings.csv` se encuentra en esta carpeta. Si se requiere
  descargarlo de nuevo, usar la versión "ml-latest-small".

## Estructura del proyecto

```
.
├── lru.go          # Nodo, lista doblemente enlazada y caché LRU
├── secuencia.go    # CargarSecuencia: lee y ordena ratings.csv
├── main.go         # Simulación: corre la secuencia y calcula el hit ratio
├── lru_test.go     # Tests unitarios y benchmarks
└── ratings.csv     # Dataset (MovieLens)
```

## Cómo ejecutar

Correr la simulación (genera la tabla de hit ratio):

```bash
go run .
```

Correr las pruebas unitarias:

```bash
go test ./...
```

Correr las pruebas con cobertura:

```bash
go test -cover ./...
```

Correr los benchmarks:

```bash
go test -bench=. -benchmem
```

## Resultados — Hit ratio

| Tamaño de caché | Hit ratio |
|---|---|
| 50 | 0.0285 |
| 100 | 0.0692 |
| 500 | 0.3226 |
| 1000 | 0.5351 |

Como se esperaba, el hit ratio crece con el tamaño de la caché: al
aumentar la capacidad, se puede retener una porción más grande del
grupo de películas más consultadas. El crecimiento entre 500 y 1000 (de 0.32 a 0.54) es más
remarcado que entre 50 y 100 (de 0.03 a 0.07): con cachés pequeñas, la
mayoría de los accesos son a películas que aún no están en la caché,
mientras que con cachés más grandes se cubre una fracción significativa
de las películas más populares del dataset.

## Complejidad

- `Get` y `Put`: **O(1)** — el mapa ubica el nodo directamente, y las
  operaciones de la lista doblemente enlazada (insertar al frente, mover al
  frente, eliminar) son manipulación de punteros, sin recorridos.
- `CargarSecuencia`: **O(n log n)** — dominado por el ordenamiento
  (`sort.Slice`) de los registros por `timestamp`.
