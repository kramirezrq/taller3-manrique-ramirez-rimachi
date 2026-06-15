package main

import (
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strconv"
)

// registrorating representa una fila de ratings.csv
// el ID de la película y el momento (timestamp) en que se calificó
type registrorating struct {
	movieID   int
	timestamp int64
}

// cargarSecuencia leemos ratings.csv, ordenamos las filas por timestamp
// retornamos la secuencia de movieId en ese orden
// formato esperado de ratings.csv: userId,movieId,rating,timestamp
func cargarSecuencia(ruta string) ([]int, error) {
	archivo, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	lector := csv.NewReader(archivo)

	if _, err := lector.Read(); err != nil {
		return nil, err
	}

	var registros []registrorating

	for {
		fila, err := lector.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		movieID, err := strconv.Atoi(fila[1])
		if err != nil {
			return nil, err
		}

		timestamp, err := strconv.ParseInt(fila[3], 10, 64)
		if err != nil {
			return nil, err
		}

		registros = append(registros, registrorating{movieID: movieID, timestamp: timestamp})
	}

	sort.Slice(registros, func(i, j int) bool {
		return registros[i].timestamp < registros[j].timestamp
	})

	secuencia := make([]int, len(registros))
	for i, r := range registros {
		secuencia[i] = r.movieID
	}

	return secuencia, nil
}
