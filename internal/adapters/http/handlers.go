package http

import (
	"net/http"

	"github.com/Vladimir-Cha/calc_rest_api/internal/core/usecases"
	"github.com/Vladimir-Cha/calc_rest_api/internal/entities"
	"github.com/labstack/echo"
)

type Handlers struct {
	calc *usecases.Calculator
}

func NewHandlers(calc *usecases.Calculator) *Handlers {
	return &Handlers{calc: calc}
}

// Result - обработчик HTTP-запросов по пути /sum
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

// GetResult возвращает сохранённый результат вычислений по токену
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

// GetTotalHandler возвращает общую сумму и произведение всех вычислений
func (h *Handlers) GetTotal(c echo.Context) error {
	total, err := h.calc.GetTotal()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, total)
}

func (h *Handlers) Token(c echo.Context) error {
	token := h.calc.GenerateToken()
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
