package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	Registry = make(map[string]string)
	mutex    sync.Mutex
)

type MockServiceInfo struct {
	Port string `json:"port"`
	URL  string `json:"url"`
}

func createMockServiceHandler(w http.ResponseWriter, r *http.Request) {
    // Отправляем запрос на получение порта
    portParamResp, err := http.Get("http://127.0.0.1:9000/port")
    if err != nil {
        // Обработка ошибки при выполнении запроса
        fmt.Println("Ошибка при выполнении запроса:", err)
        http.Error(w, "Failed to fetch port parameter", http.StatusInternalServerError)
        return
    }
    defer portParamResp.Body.Close()

    // Проверяем статус ответа
    if portParamResp.StatusCode != http.StatusOK {
        // Обработка ошибки статуса ответа
        fmt.Println("Ошибка статуса ответа:", portParamResp.Status)
        http.Error(w, "Failed to fetch port parameter: "+portParamResp.Status, portParamResp.StatusCode)
        return
    }

    // Выводим содержимое тела ответа в консоль
    bodyContent, err := ioutil.ReadAll(portParamResp.Body)
    if err != nil {
        // Обработка ошибки чтения тела ответа
        fmt.Println("Ошибка при чтении тела ответа:", err)
        http.Error(w, "Failed to read port parameter", http.StatusInternalServerError)
        return
    }

    // Извлекаем порт из строки ответа
    portParamStr := strings.TrimSpace(string(bodyContent))
    parts := strings.Split(portParamStr, " ")
    if len(parts) != 3 {
        // Неправильный формат ответа
        fmt.Println("Неправильный формат ответа:", portParamStr)
        http.Error(w, "Invalid response format", http.StatusInternalServerError)
        return
    }

    port := parts[2]

    // Проверяем, что параметр порта не пустой
    if port == "" {
        http.Error(w, "Port parameter is missing", http.StatusBadRequest)
        return
    }

    // Преобразуем строку порта в число
    mockPort, err := strconv.Atoi(port)
    if err != nil {
        http.Error(w, "Invalid port parameter", http.StatusBadRequest)
        return
    }

	cmd := exec.Command("go", "run", "mock/mock_service.go", "--port", strconv.Itoa(mockPort))
	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to create mock service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	Registry[strconv.Itoa(mockPort)] = fmt.Sprintf("http://localhost:%d", mockPort)

	fmt.Printf("Mock service created successfully! Port: %d\n", mockPort)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Mock service created successfully"))
}

func getMockServicesInfoHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var mockServicesInfo []MockServiceInfo

	for port, url := range Registry {
		mockServicesInfo = append(mockServicesInfo, MockServiceInfo{
			Port: port,
			URL:  url,
		})
	}

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
