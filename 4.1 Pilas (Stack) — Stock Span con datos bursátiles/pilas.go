package pilas

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	top := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return top, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func CalcularStockSpan(precios []float64) []int {
	n := len(precios)
	span := make([]int, n)
	var st Stack[int]

	for i := 0; i < n; i++ {
		for top, ok := st.Peek(); ok && precios[top] <= precios[i]; top, ok = st.Peek() {
			st.Pop()
		}
		if st.IsEmpty() {
			span[i] = i + 1
		} else {
			top, _ := st.Peek()
			span[i] = i - top
		}
		st.Push(i)
	}
	return span
}

type Registro struct {
	Fecha  string
	Precio float64
}

func LeerPrecios(ruta string) ([]float64, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, fmt.Errorf("LeerPrecios: no se pudo abrir el archivo: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("LeerPrecios: error leyendo encabezado: %w", err)
	}

	closeIdx := -1
	for i, col := range header {
		if col == "Close" {
			closeIdx = i
			break
		}
	}
	if closeIdx == -1 {
		return nil, errors.New("LeerPrecios: columna 'Close' no encontrada")
	}

	var precios []float64
	for {
		row, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("LeerPrecios: error leyendo fila: %w", err)
		}
		precio, err := strconv.ParseFloat(row[closeIdx], 64)
		if err != nil {
			continue
		}
		precios = append(precios, precio)
	}

	if len(precios) == 0 {
		return nil, errors.New("LeerPrecios: el archivo no contiene precios válidos")
	}
	return precios, nil
}

func LeerRegistros(ruta string) ([]Registro, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, fmt.Errorf("LeerRegistros: no se pudo abrir el archivo: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("LeerRegistros: error leyendo encabezado: %w", err)
	}

	dateIdx, closeIdx := -1, -1
	for i, col := range header {
		switch col {
		case "Date":
			dateIdx = i
		case "Close":
			closeIdx = i
		}
	}
	if closeIdx == -1 {
		return nil, errors.New("LeerRegistros: columna 'Close' no encontrada")
	}

	var registros []Registro
	for {
		row, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("LeerRegistros: error leyendo fila: %w", err)
		}
		precio, err := strconv.ParseFloat(row[closeIdx], 64)
		if err != nil {
			continue
		}
		reg := Registro{Precio: precio}
		if dateIdx >= 0 && dateIdx < len(row) {
			reg.Fecha = row[dateIdx]
		}
		registros = append(registros, reg)
	}

	if len(registros) == 0 {
		return nil, errors.New("LeerRegistros: el archivo no contiene registros válidos")
	}
	return registros, nil
}

func DiaMayorSpan(spans []int) int {
	if len(spans) == 0 {
		return -1
	}
	maxIdx := 0
	for i := 1; i < len(spans); i++ {
		if spans[i] > spans[maxIdx] {
			maxIdx = i
		}
	}
	return maxIdx
}
