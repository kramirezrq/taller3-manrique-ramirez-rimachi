package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Pelicula struct {
	ID     int
	Titulo string
	Prom   float64 // Estandarizado a 'Prom'
}

type NodoAVL struct {
	clave float64
	datos []Pelicula // Estandarizado a 'datos'
	alt   int
	izq   *NodoAVL
	der   *NodoAVL
}

func Altura(n *NodoAVL) int {
	if n == nil {
		return -1
	}
	return n.alt
}

func factorBalance(n *NodoAVL) int {
	return Altura(n.izq) - Altura(n.der)
}

func actualizarAltura(n *NodoAVL) {
	izq := Altura(n.izq)
	der := Altura(n.der)
	if izq > der {
		n.alt = 1 + izq
	} else {
		n.alt = 1 + der
	}
}

func rotDer(y *NodoAVL) *NodoAVL {
	X := y.izq
	k2 := X.der
	X.der = y
	y.izq = k2
	actualizarAltura(y)
	actualizarAltura(X) // Corregido 'x' por 'X'
	return X
}

func rotIzq(x *NodoAVL) *NodoAVL {
	Y := x.der
	k2 := Y.izq
	Y.izq = x
	x.der = k2
	actualizarAltura(x)
	actualizarAltura(Y)
	return Y
}

func balancear(n *NodoAVL) *NodoAVL {
	actualizarAltura(n)
	fb := factorBalance(n)

	if fb > 1 && factorBalance(n.izq) >= 0 {
		return rotDer(n)
	}
	if fb > 1 && factorBalance(n.izq) < 0 {
		n.izq = rotIzq(n.izq)
		return rotDer(n)
	}
	if fb < -1 && factorBalance(n.der) <= 0 {
		return rotIzq(n)
	}
	if fb < -1 && factorBalance(n.der) > 0 {
		n.der = rotDer(n.der)
		return rotIzq(n)
	}
	return n
}

func Insertar(raiz *NodoAVL, clave float64, dato Pelicula) *NodoAVL {
	if raiz == nil {
		return &NodoAVL{
			clave: clave,
			datos: []Pelicula{dato},
			alt:   0,
		}
	}

	if clave < raiz.clave {
		raiz.izq = Insertar(raiz.izq, clave, dato)
	} else if clave > raiz.clave {
		raiz.der = Insertar(raiz.der, clave, dato)
	} else {
		raiz.datos = append(raiz.datos, dato)
	}
	return balancear(raiz)
}

func ConsultaRango(raiz *NodoAVL, a, b float64) []Pelicula {
	if raiz == nil {
		return []Pelicula{}
	}

	resultado := []Pelicula{}

	if raiz.clave > a {
		resultado = append(resultado, ConsultaRango(raiz.izq, a, b)...)
	}

	if raiz.clave >= a && raiz.clave <= b {
		resultado = append(resultado, raiz.datos...)
	}

	if raiz.clave < b {
		resultado = append(resultado, ConsultaRango(raiz.der, a, b)...)
	}

	return resultado
}

func main() {
	var rutaItem, rutaData string

	fmt.Print("Introducir ruta del archivo u.item: ")
	fmt.Scanln(&rutaItem)
	rutaItem = strings.Trim(rutaItem, `"'`)
	archivo, err := os.Open(rutaItem)
	if err != nil {
		fmt.Println("Error al abrir u.item:", err)
		return
	}
	defer archivo.Close()

	fmt.Print("Introduce la ruta del archivo u.data: ")
	fmt.Scanln(&rutaData)
	rutaData = strings.Trim(rutaData, `"'`)
	archivo2, err := os.Open(rutaData)
	if err != nil {
		fmt.Println("Error al abrir u.data:", err)
		return
	}
	defer archivo2.Close()

	fmt.Println("Archivos abiertos correctamente")

	scanner := bufio.NewScanner(archivo)
	peliculas := make(map[int]string)
	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Split(linea, "|")
		if len(partes) < 2 {
			continue
		}
		ID, err := strconv.Atoi(partes[0])
		if err != nil {
			continue
		}
		peliculas[ID] = partes[1]
	}

	fmt.Printf("Se cargaron %d películas correctamente.\n", len(peliculas))

	scanner2 := bufio.NewScanner(archivo2)
	sumas := make(map[int]float64)
	conteos := make(map[int]int)
	for scanner2.Scan() {
		linea := scanner2.Text()
		partes := strings.Split(linea, "\t")
		if len(partes) < 3 { // Validación de seguridad añadida
			continue
		}
		movieID, err := strconv.Atoi(partes[1])
		if err != nil {
			continue
		}
		rating, err := strconv.ParseFloat(partes[2], 64)
		if err != nil {
			continue
		}
		sumas[movieID] += rating
		conteos[movieID]++
	}

	promedios := make(map[int]float64)
	for ID, suma := range sumas {
		promedios[ID] = suma / float64(conteos[ID])
	}

	lstPeli := []Pelicula{}
	for ID, titulo := range peliculas {
		lstPeli = append(lstPeli, Pelicula{
			ID:     ID,
			Titulo: titulo,
			Prom:   promedios[ID], // Corregido
		})
	}
	fmt.Println("Totales con P:", len(lstPeli))

	var RaizAVL *NodoAVL
	for _, pelicula := range lstPeli {
		RaizAVL = Insertar(RaizAVL, pelicula.Prom, pelicula) // Corregido pelicula.P -> pelicula.Prom
	}

	fmt.Println("AVL construido. Altura:", Altura(RaizAVL))

	fmt.Println("\nPeliculas con un rating entre 4.5 y 5.0:")
	resultados := ConsultaRango(RaizAVL, 4.5, 5.0)
	fmt.Println("Total:", len(resultados))
	for _, p := range resultados {
		fmt.Printf("  %s - %.2f\n", p.Titulo, p.Prom) // Corregido p.P -> p.Prom
	}
}
