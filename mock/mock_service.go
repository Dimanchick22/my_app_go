package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
    // Определяем флаг для порта
    portPtr := flag.Int("port", 8000, "Port number for the mock service")
    flag.Parse()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from mock service on port %d!", *portPtr)
    })

    // Преобразуем порт в строку для передачи в ListenAndServe
    port := ":" + strconv.Itoa(*portPtr)
    fmt.Printf("Mock service is running on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
