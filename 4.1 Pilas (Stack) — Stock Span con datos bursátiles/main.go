package main

import (
	"fmt"
	"os"
	"strconv"
	"taller3/pilas"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Uso: go run ./pilas/main <ruta_dataset> [N]")
		os.Exit(1)
	}
	ruta := os.Args[1]
	n := 5
	if len(os.Args) >= 3 {
		if v, err := strconv.Atoi(os.Args[2]); err == nil && v > 0 {
			n = v
		}
	}

	registros, err := pilas.LeerRegistros(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	precios := make([]float64, len(registros))
	for i, r := range registros {
		precios[i] = r.Precio
	}
	spans := pilas.CalcularStockSpan(precios)

	fmt.Printf("Dataset: %s | Total días: %d\n\n", ruta, len(registros))
	fmt.Printf("%-12s  %-10s  %s\n", "Fecha", "Cierre", "Span")
	fmt.Println("─────────────────────────────────")

	limite := n
	if limite > len(registros) {
		limite = len(registros)
	}
	fmt.Printf("── Primeros %d días ──\n", limite)
	for i := 0; i < limite; i++ {
		fmt.Printf("%-12s  %10.4f  %d\n", registros[i].Fecha, registros[i].Precio, spans[i])
	}
	if len(registros) > 2*limite {
		fmt.Printf("\n── Últimos %d días ──\n", limite)
		for i := len(registros) - limite; i < len(registros); i++ {
			fmt.Printf("%-12s  %10.4f  %d\n", registros[i].Fecha, registros[i].Precio, spans[i])
		}
	}

	maxIdx := pilas.DiaMayorSpan(spans)
	fmt.Printf("\nDía con mayor span: %s  Cierre=%.4f  Span=%d\n",
		registros[maxIdx].Fecha, registros[maxIdx].Precio, spans[maxIdx])
}
