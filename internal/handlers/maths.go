package handlers

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// Хранилище данных
var (
	storage = struct {
		mu      sync.Mutex
		results map[string]NumResponse
	}{
		results: make(map[string]NumResponse),
	}
)

var (
	total = struct {
		mu  sync.Mutex
		sum float64
		mul float64
	}{
		sum: 0,
		mul: 0,
	}
)

// Result - обработчик HTTP-запросов по пути /sum
func Result(c echo.Context) error {
	token := c.Request().Header.Get("Token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, "Нет токена")
	}

	var req NumRequest
	//декодер для тела запроса
	if err := c.Bind(&req); err != nil { //парсим JSON в структуру SumOfNumbers
		return c.JSON(http.StatusBadRequest, "Неверный формат JSON"+err.Error())
	}

	//Считаем сумму и произведение чисел
	var sum float64
	var mul float64
	sum, mul = 0.0, 1.0

	for _, num := range req.Numbers {
		sum += num
		mul *= num
	}

	//сохраняем результат
	storage.mu.Lock()
	storage.results[token] = NumResponse{
		ResponseNumbers:       sum,
		MultiplicationNumbers: mul,
	}
	storage.mu.Unlock()

	//обновляем общие суммы
	total.mu.Lock()
	total.sum += sum
	total.mul += mul
	total.mu.Unlock()

	//Формируем ответ
	//response := map[string]float64{"sum": sum} //мапа для ответа
	return c.JSON(http.StatusOK, NumResponse{
		ResponseNumbers:       sum,
		MultiplicationNumbers: mul,
	})

}

// GetResultHandler возвращает сохранённый результат вычислений по токену
func GetResult(c echo.Context) error {
	//Условие для проверки написания token и Token
	token := c.QueryParam("token")
	if token == "" {
		token = c.QueryParam("Token")
	}

	if token == "" {
		return c.JSON(http.StatusBadRequest, "Нет токена")
	}

	// Результат вычислений для токена
	storage.mu.Lock()
	result, exist := storage.results[token]
	storage.mu.Unlock()
	if !exist {
		return c.JSON(http.StatusNotFound, "Нет данных для токена")
	}

	//Результат общих вычислений для токена
	total.mu.Lock()
	totalResult := TotalResult{
		TotalSum: total.sum,
		TotalMul: total.mul,
	}
	total.mu.Unlock()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"individual": result,
		"total":      totalResult,
	})
}

// GetTotalHandler возвращает общую сумму и произведение всех вычислений
func GetTotal(c echo.Context) error {

	total.mu.Lock()
	defer total.mu.Unlock()

	return c.JSON(http.StatusOK, TotalResult{
		TotalSum: total.sum,
		TotalMul: total.mul,
	})
}

// Получение токена
func Token(c echo.Context) error {
	token := uuid.New().String()
	return c.JSON(http.StatusOK, TokenResponse{Token: token})
}
