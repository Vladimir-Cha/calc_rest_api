package handlers

import (
	"net/http"

	"github.com/Vladimir-Cha/calc_rest_api/internal/core/usecases"
	"github.com/Vladimir-Cha/calc_rest_api/internal/entities"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	calc *usecases.Calculator
}

func NewHandlers(calc *usecases.Calculator) *Handlers {
	return &Handlers{calc: calc}
}

// @Summary Выполнить вычисления
// @Description Принимает массив чисел и возвращает сумму и произведение
// @Tags calculations
// @Accept  json
// @Produce json
// @Param   numbers body     entities.NumRequest true "Массив чисел"
// @Param   Token   header   string              true "Токен доступа"
// Result - обработчик HTTP-запросов по пути /result
// @Success 200 {object} entities.NumResponse
// @Failure 400 {string} string "Неверный формат JSON"
// @Failure 400 {string} string "Нет токена"
// @Router /result [post]
func (h *Handlers) Result(c echo.Context) error {
	token := c.Request().Header.Get("Token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, "Нет токена")
	}

	var req entities.NumRequest
	//декодер для тела запроса
	if err := c.Bind(&req); err != nil { //парсим JSON в структуру SumOfNumbers
		return c.JSON(http.StatusBadRequest, "Неверный формат JSON"+err.Error())
	}

	res, err := h.calc.Result(token, req.Numbers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// @Summary Получить общие результаты вычислений по токену
// @Description Возвращает индивидуальные и общие значения всех выполненных операций сумм и произведений по токену
// @Tags calculations
// @Produce json
// @Param   token  query   string  true  "Токен доступа"
// Result - обработчик HTTP-запросов по пути /tokenresult
// @Success 200 {object} map[string]interface{} "Пример: {"individual": {"sum": 10, "multipl": 24}, "total": {"totalsum": 100, "totalmultipl": 1000}}"
// @Failure 400 {string} string "Нет данных для токена"
// @Router /tokenresult [get]
func (h *Handlers) GetResult(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		token = c.QueryParam("Token") // Проверяем 2 варианта написания
	}

	res, total, err := h.calc.GetResult(token)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Нет данных для токена")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"individual": res,
		"total":      total,
	})
}

// @Summary Получить общие результаты вычислений
// @Description Возвращает общие значения всех выполненных операций сумм и произведений
// @Tags calculations
// @Produce json
// Result - обработчик HTTP-запросов по пути /totalresult
// @Success 200 {object} entities.TotalResult "Пример: {"TotalSum": 100, "TotalMul": 1000}"
// @Failure 500 {string} string "{"error": "Описание ошибки"}"
// @Router /totalresult [get]
func (h *Handlers) GetTotal(c echo.Context) error {
	total, err := h.calc.GetTotal()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, total)
}

// @Summary Генерация токена
// @Description Создает новый уникальный токен для доступа к API
// @Tags auth
// @Accept  json
// @Produce json
// @Success 200 {object} map[string]string "Пример: {"token": "01735715-8853-48d8-9c7e-a43e60ca90ef"}"
// @Router /token [post]
func (h *Handlers) Token(c echo.Context) error {
	token := h.calc.GenerateToken()
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
