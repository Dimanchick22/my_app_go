package mock

import (
	"fmt"
	"net/http"
)

// HandleMockRequest обрабатывает запросы к мок-сервису
func HandleMockRequest(w http.ResponseWriter, r *http.Request) {
    // Ваша логика обработки запросов к мок-сервису
    fmt.Fprintf(w, "Hello from mock service!")
}
