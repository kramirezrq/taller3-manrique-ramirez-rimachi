# Rate Limiter para Logs Apache en Go

## Descripción

Este proyecto implementa un **Rate Limiter** basado en ventanas de tiempo deslizantes (*Sliding Window*) para analizar archivos de logs Apache.

El programa procesa un archivo de registros y determina si una petición debe ser aceptada o rechazada según las siguientes reglas:

* Cada dirección IP puede realizar como máximo **M peticiones** dentro de un intervalo de **T segundos**.
* Si una IP supera ese límite, la petición es rechazada.
* Al finalizar, el programa genera estadísticas sobre el procesamiento realizado.

---

## Características

* Procesamiento eficiente de archivos de logs de gran tamaño.
* Implementación de colas para gestionar las ventanas de tiempo por IP.
* Eliminación automática de registros fuera de la ventana de análisis.
* Conteo global de peticiones rechazadas.
* Ranking de las 5 IPs con más rechazos.
* Medición del tiempo total de ejecución.

---

## Estructura del Proyecto

```text
.
├── main.go
└── README.md
```

---

## Formato del Log

El programa espera registros con formato Apache similar al siguiente:

```text
192.168.1.10 - - [10/Oct/2025:13:55:36 -0500] "GET /index.html HTTP/1.1" 200 1024
```

De cada línea se extraen:

* Dirección IP
* Fecha y hora de la petición

---

## Algoritmo Utilizado

### Sliding Window

Para cada IP se mantiene una cola de timestamps.

Proceso:

1. Se elimina de la cola cualquier petición cuya antigüedad sea mayor que `T` segundos.
2. Se verifica cuántas peticiones permanecen dentro de la ventana.
3. Si el número es mayor o igual a `M`, la petición es rechazada.
4. En caso contrario, se registra la nueva petición y se acepta.

Complejidad aproximada:

* Tiempo: **O(n)**
* Memoria: **O(k)**

Donde:

* `n` = número total de peticiones.
* `k` = número de peticiones activas dentro de las ventanas temporales.

---

## Requisitos

* Go 1.18 o superior.

Verificar instalación:

```bash
go version
```

---

## Compilación

```bash
go build main.go
```

Generará un ejecutable:

```bash
main
```

En Windows:

```bash
main.exe
```

---

## Ejecución

```bash
go run main.go <archivo_log> <M> <T>
```

o usando el ejecutable:

```bash
./main <archivo_log> <M> <T>
```

### Parámetros

| Parámetro   | Descripción                            |
| ----------- | -------------------------------------- |
| archivo_log | Ruta del archivo de logs Apache        |
| M           | Máximo número de peticiones permitidas |
| T           | Ventana de tiempo en segundos          |

---

## Ejemplo de Uso

```bash
go run main.go access.log 100 60
```

Interpretación:

* Máximo 100 peticiones
* Dentro de una ventana de 60 segundos
* Por cada dirección IP

---

## Salida Esperada

```text
PROCESANDO PETICIONES

RESUMEN FINAL

Tiempo total: 1.23s
Total de peticiones analizadas: 250000
Total global de rechazos: 1234

Top 5 IPs con más rechazos:

1. IP: 192.168.1.10 | Rechazos: 500
2. IP: 192.168.1.15 | Rechazos: 320
3. IP: 10.0.0.5     | Rechazos: 180
4. IP: 172.16.1.8   | Rechazos: 140
5. IP: 8.8.8.8      | Rechazos: 94
```

---

## Componentes Principales

### Cola

Estructura utilizada para almacenar los timestamps de cada IP.

Métodos:

* `Enqueue()`
* `Dequeue()`
* `Front()`
* `Len()`

### PermitirPeticion()

Determina si una petición puede procesarse o debe ser rechazada según la política de Rate Limiting.

### ParsearLinea()

Extrae:

* Dirección IP
* Timestamp Unix

a partir de una línea del log Apache.

---

## Casos de Error Controlados

* Cantidad incorrecta de argumentos.
* Valores no enteros para `M` o `T`.
* Archivo inexistente o inaccesible.
* Líneas de log mal formateadas.
* Errores de lectura durante el procesamiento.

---

## Autor

Proyecto desarrollado en Go para el análisis de tráfico y control de peticiones mediante un algoritmo de Rate Limiting basado en ventanas deslizantes.

