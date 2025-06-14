package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Vladimir-Cha/calc_rest_api/internal/adapters/handlers"
	"github.com/Vladimir-Cha/calc_rest_api/internal/adapters/storage"
	"github.com/Vladimir-Cha/calc_rest_api/internal/core/usecases"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// getURL формирует URL для тестовых запросов
func getURL(path string) string {
	port := 8080
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if eport, err := strconv.Atoi(envPort); err == nil {
			port = eport
		}
	}
	path = strings.TrimPrefix(strings.ReplaceAll(path, `\`, `/`), "/") //меняем слэши, если неправильные
	return fmt.Sprintf("http://localhost:%d/%s", port, path)           //формирование URL
}

func TestApp(t *testing.T) {
	// Запуск сервера в отдельной горутине с обработчиками сервера
	go func() {
		storage := storage.NewMathStorage()
		calculator := usecases.NewCalculator(storage)
		h := handlers.NewHandlers(calculator)

		e := echo.New()

		// тестовый обработчик для корневого запроса
		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Server is running")
		})

		// Регистрация маршрутов
		e.POST("/result", h.Result)
		e.GET("/totalres", h.GetResult)
		e.POST("/token", h.Token)

		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Даем серверу время на запуск
	time.Sleep(500 * time.Millisecond)

	// Тест доступности сервера
	t.Run("Server availability", func(t *testing.T) {
		resp, err := http.Get(getURL(""))
		if !assert.NoError(t, err, "Server should be available") {
			return
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200")
	})

	// Тест генерации токена
	t.Run("Token generation", func(t *testing.T) {
		resp, err := http.Post(getURL("/token"), "application/json", nil)
		if !assert.NoError(t, err, "Token request failed") {
			return
		}
		defer resp.Body.Close()

		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		if assert.NoError(t, err, "Should decode JSON response") {
			assert.NotEmpty(t, result["token"], "Token should not be empty")
		}
	})
}
