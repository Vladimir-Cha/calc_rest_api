package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func MultimplHandler(c echo.Context) error {
	var req SumOfNumbers
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Неверный формат JSON"+err.Error())
	}

	var mul float64
	mul = 1
	for _, num := range req.Numbers {
		mul *= num
	}
	mu.Lock()
	totalMul += mul
	mu.Unlock()
	return c.JSON(http.StatusOK, ResponseSumOfNumbers{MultiplicationNumbers: mul})
}

func GetMulHandler(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()
	return c.JSON(http.StatusOK, ResponseSumOfNumbers{MultiplicationNumbers: totalMul})
}
