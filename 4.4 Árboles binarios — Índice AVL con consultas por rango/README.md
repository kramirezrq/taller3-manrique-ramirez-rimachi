# Árbol AVL para Gestión de Películas

## Descripción

Programa desarrollado en Go que implementa un Árbol AVL para almacenar y consultar películas según su calificación promedio. Los datos son obtenidos de los archivos `u.item` y `u.data` del conjunto de datos MovieLens.

## Funcionalidades

* Carga de películas y calificaciones.
* Cálculo del rating promedio de cada película.
* Inserción de películas en un Árbol AVL.
* Balanceo automático mediante rotaciones AVL.
* Consulta eficiente de películas por rango de calificaciones.

## Estructuras Principales

* **Pelicula:** almacena ID, título y promedio.
* **NodoAVL:** almacena la clave (promedio), las películas asociadas y enlaces a los nodos hijos.

## Ejecución

Compilar y ejecutar:

```bash
go run arboles.go
```

Luego ingresar las rutas de los archivos:

```text
u.item
u.data
```

## Complejidad

| Operación          | Complejidad  |
| ------------------ | ------------ |
| Inserción          | O(log n)     |
| Búsqueda por rango | O(log n + k) |

Donde `n` es el número de películas almacenadas y `k` la cantidad de resultados encontrados.

## Autor

Proyecto desarrollado como práctica de implementación de árboles AVL en Go.
