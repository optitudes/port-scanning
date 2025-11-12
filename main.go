package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		fmt.Println("Scanning port ",p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		fmt.Println("Closing connection with port ",p)
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int
	var wg sync.WaitGroup

	// 1. Iniciar workers (100 goroutines concurrentes)
	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ports, results)
		}()
	}

	// 2. Goroutine para enviar puertos a escanear
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i // Aquí ENVÍAS los puertos
		}
		close(ports) // Cerrar cuando termines de enviar
	}()

	// 3. Goroutine para cerrar results cuando todos los workers terminen
	go func() {
		wg.Wait()      // Esperar a que todos los workers terminen
		close(results) // Entonces cerrar results
	}()

	// 4. Recibir resultados (esto bloquea hasta que results se cierre)
	for port := range results {
		if port != 0 {
			openports = append(openports, port)
		}
	}

	// 5. Mostrar resultados
	sort.Ints(openports)
	fmt.Printf("\nPuertos abiertos en scanme.nmap.org:\n")
	fmt.Println("=====================================")
	for _, port := range openports {
		fmt.Printf("Puerto %d está ABIERTO\n", port)
	}
	fmt.Printf("\n✓ Total: %d puertos abiertos de 1024 escaneados\n", len(openports))
}
