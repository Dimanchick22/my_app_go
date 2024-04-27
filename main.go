package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"sync"
)

var (
    Registry      = make(map[string]string)
    mutex         sync.Mutex
    mockPortBase  = 8000 // Базовый порт для мок-сервисов
    mockServices  int    // Переменная для отслеживания количества созданных мок-сервисов
    mockServicesMutex sync.Mutex // Mutex для безопасного доступа к mockServices
)

type MockServiceInfo struct {
    Port string `json:"port"`
    URL  string `json:"url"`
}

func createMockServiceHandler(w http.ResponseWriter, r *http.Request) {
    // Увеличиваем счетчик мок-сервисов и получаем порт для нового мок-сервиса
    mockServicesMutex.Lock()
    mockPort := mockPortBase + mockServices
    mockServices++
    numMockServices := mockServices
    mockServicesMutex.Unlock()

    // Создаем команду для запуска мок-сервиса с новым портом
    cmd := exec.Command("go", "run", "mock/mock_service.go", "--port", strconv.Itoa(mockPort))

    // Запускаем команду
    if err := cmd.Start(); err != nil {
        http.Error(w, "Failed to create mock service: "+err.Error(), http.StatusInternalServerError)
        return
    }

    mutex.Lock()
    defer mutex.Unlock()
    Registry[strconv.Itoa(mockPort)] = fmt.Sprintf("http://localhost:%d", mockPort)

    // Выводим информацию о создании мок-сервиса
    fmt.Printf("Mock service created successfully! Port: %d, Total mock services: %d\n", mockPort, numMockServices)

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Mock service created successfully"))
}

func getMockServicesInfoHandler(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    defer mutex.Unlock()

    // Создаем слайс для хранения информации о мок-сервисах
    var mockServicesInfo []MockServiceInfo

    // Итерируемся по реестру мок-сервисов и заполняем информацию
    for port, url := range Registry {
        mockServicesInfo = append(mockServicesInfo, MockServiceInfo{
            Port: port,
            URL:  url,
        })
    }

    // Кодируем информацию в JSON и отправляем клиенту
    jsonBytes, err := json.Marshal(mockServicesInfo)
    if err != nil {
        http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}


func main() {
    http.HandleFunc("/createMockService", createMockServiceHandler)
    http.HandleFunc("/getMockServicesInfo", getMockServicesInfoHandler)

    log.Fatal(http.ListenAndServe(":7000", nil))
}
