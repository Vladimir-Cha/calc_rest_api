package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	path = strings.ReplaceAll(path, `\`, `/`) //меняем слэши, если неправильные
	path = strings.TrimPrefix(path, "/")      //убираем "/" из конца ссылки, если есть
	//формирование URL
	u := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d", port),
		Path:   path,
	}
	return u.String()
}

func startAPP() *echo.Echo {
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

	go func() {
		e.Start(":8080")
	}()

	time.Sleep(1 * time.Second)
	return e
}
func TestApp(t *testing.T) {
	startAPP()

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
	// Тест вычислений
	t.Run("Calculation", func(t *testing.T) {
		// Сначала получаем токен
		tokenResp, err := http.Post(getURL("/token"), "application/json", nil)
		if !assert.NoError(t, err) {
			return
		}
		defer tokenResp.Body.Close()

		var tokenData map[string]string
		err = json.NewDecoder(tokenResp.Body).Decode(&tokenData)
		if !assert.NoError(t, err) {
			return
		}

		// Отправляем данные для расчета
		requestData := map[string][]float64{"numbers": {1, 2, 3, 4}}
		body, _ := json.Marshal(requestData)

		req, err := http.NewRequest("POST", getURL("/result"), bytes.NewBuffer(body))
		if !assert.NoError(t, err) {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Token", tokenData["token"])

		resp, err := http.DefaultClient.Do(req)
		if !assert.NoError(t, err) {
			return
		}
		defer resp.Body.Close()

		var result map[string]float64
		err = json.NewDecoder(resp.Body).Decode(&result)
		if assert.NoError(t, err, "Should decode JSON response") {
			assert.Equal(t, 10.0, result["sum"], "Sum should be 10 (1+2+3+4)")
			assert.Equal(t, 24.0, result["multipl"], "Multiplication should be 24 (1*2*3*4)")
		}
	})
}
