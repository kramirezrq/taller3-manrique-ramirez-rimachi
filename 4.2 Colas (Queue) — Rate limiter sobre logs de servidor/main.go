package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Cola struct {
	items []int64
}

func (c *Cola) Enqueue(ts int64) {
	c.items = append(c.items, ts)
}

func (c *Cola) Dequeue() (int64, bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	valor := c.items[0]
	c.items = c.items[1:]
	return valor, true
}

func (c *Cola) Front() (int64, bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	return c.items[0], true
}

func (c *Cola) Len() int {
	return len(c.items)
}

func PermitirPeticion(colas map[string]*Cola, ip string, ts int64, M int, T int64) bool {
	cola, existe := colas[ip]
	if !existe {
		cola = &Cola{}
		colas[ip] = cola
	}

	for {
		frontTs, tieneElementos := cola.Front()
		if !tieneElementos || frontTs > (ts-T) {
			break
		}
		cola.Dequeue()
	}

	if cola.Len() >= M {
		return false
	}

	cola.Enqueue(ts)
	return true
}

func ParsearLinea(linea string) (string, int64, error) {
	partes := strings.Split(linea, " ")
	if len(partes) < 4 {
		return "", 0, fmt.Errorf("línea mal formateada")
	}

	ip := partes[0]
	fechaStr := partes[3]
	fechaStr = strings.TrimPrefix(fechaStr, "[")
	zonaStr := strings.TrimSuffix(partes[4], "]")
	fechaCompleta := fechaStr + " " + zonaStr
	layoutApache := "02/Jan/2006:15:04:05 -0700"
	t, err := time.Parse(layoutApache, fechaCompleta)
	if err != nil {
		return "", 0, err
	}

	return ip, t.Unix(), nil
}

type IpRechazos struct {
	IP       string
	Rechazos int
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Uso:")
		fmt.Println("go run main.go <archivo_log> <M> <T>")
		return
	}
	rutaArchivo := os.Args[1]
	M, errM := strconv.Atoi(os.Args[2])
	T, errT := strconv.ParseInt(os.Args[3], 10, 64)

	if errM != nil || errT != nil {
		fmt.Println("Error: M y T deben ser enteros.")
		return
	}

	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		fmt.Printf("Error'%s': %v\n", rutaArchivo, err)
		return
	}
	defer archivo.Close()

	colas := make(map[string]*Cola)
	rechazosPorIP := make(map[string]int)

	totalLineas := 0
	totalRechazos := 0
	tiempoInicio := time.Now()

	scanner := bufio.NewScanner(archivo)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	fmt.Println(" PROCESANDO PETICIONES ")

	for scanner.Scan() {
		totalLineas++
		linea := scanner.Text()

		ip, ts, err := ParsearLinea(linea)
		if err != nil {
			continue
		}

		if !PermitirPeticion(colas, ip, ts, M, T) {
			totalRechazos++
			rechazosPorIP[ip]++
		}

		if totalLineas%500000 == 0 {
			fmt.Printf("\rLíneas procesadas: %d...", totalLineas)
		}
	}

	var listaRechazos []IpRechazos
	for ip, count := range rechazosPorIP {
		listaRechazos = append(listaRechazos, IpRechazos{IP: ip, Rechazos: count})
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\nError durante la lectura: %v\n", err)
	}

	sort.Slice(listaRechazos, func(i, j int) bool {
		return listaRechazos[i].Rechazos > listaRechazos[j].Rechazos
	})
	duracion := time.Since(tiempoInicio)
	fmt.Println("\n RESUMEN FINAL ")
	fmt.Printf("Tiempo total:               %v\n", duracion)
	fmt.Printf("Total de peticiones analizadas: %d\n", totalLineas)
	fmt.Printf("Total global de rechazos:       %d\n", totalRechazos)
	fmt.Println("\nTop 5 IPs con más rechazos:")

	limiteTop := 5
	if len(listaRechazos) < limiteTop {
		limiteTop = len(listaRechazos)
	}

	if limiteTop == 0 {
		fmt.Println("No se registraron rechazos.")
	} else {
		for i := 0; i < limiteTop; i++ {
			fmt.Printf("%d. IP: %-15s | Rechazos: %d\n", i+1, listaRechazos[i].IP, listaRechazos[i].Rechazos)
		}
	}
}
