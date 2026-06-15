package main

import "fmt"

func main() {
	rutaCSV := "ratings.csv"

	secuencia, err := cargarSecuencia(rutaCSV)
	if err != nil {
		fmt.Println("error al cargar la secuencia:", err)
		return
	}

	fmt.Printf("secuencia cargada: %d accesos\n\n", len(secuencia))

	capacidades := []int{50, 100, 500, 1000}

	fmt.Println("tamaño caché | hit ratio")

	for _, capacidad := range capacidades { //iteramos sobre distintas capacidades
		cache, err := NuevaLRU(capacidad) //creamos cache con nueva capacidad
		if err != nil {
			fmt.Println("error al crear la caché:", err)
			return
		}

		hits := 0
		for _, movieID := range secuencia {
			_, encontrado := cache.Get(movieID) //preguntamos si la película ya está en caché
			if encontrado {
				hits++ //fue un hit, sumamos 1
			} else {
				cache.Put(movieID, movieID) //fue un miss, la guardamos con Put
			}
		}

		hitRatio := float64(hits) / float64(len(secuencia)) //hits dice en que medida es util la cache
		fmt.Printf("%12d | %.4f\n", capacidad, hitRatio)
	}
}
