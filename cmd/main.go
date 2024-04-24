package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
    "github.com/dimanchick22/GO/pkg/mock"
)

// Registry хранит информацию о зарегистрированных мок-сервисах
var Registry = make(map[string]string)
var mutex sync.Mutex

func createMockServiceHandler(w http.ResponseWriter, r *http.Request) {
    // Парсим параметры запроса для создания мок-сервиса
    r.ParseForm()
    mockPath := r.Form.Get("path")
    mockPort := r.Form.Get("port")

    // Создаем новый мок-сервис
    cmd := exec.Command("go", "run", "mock_service.go") // Выполнение мок-сервиса
    if err := cmd.Start(); err != nil {
        http.Error(w, "Failed to create mock service: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Регистрируем новый мок-сервис в реестре
    mutex.Lock()
    defer mutex.Unlock()
    Registry[mockPath] = fmt.Sprintf("http://localhost:%s", mockPort)

    // Возвращаем успешный ответ
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Mock service created successfully"))
}

func main() {
    http.HandleFunc("/createMockService", createMockServiceHandler)
    http.HandleFunc("/mock", mock.HandleMockRequest) // Обработка запросов к мок-сервису

    log.Fatal(http.ListenAndServe(":8080", nil))
}
