package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

// sumHandler - обработчик HTTP-запросов по пути /sum
func SumHandler(c echo.Context) error {
	//Если метод не POST, то выводим ошибку
	/*if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, "Метод должен быть POST")
	}
	*/
	var req SumOfNumbers
	//декодер для тела запроса
	if err := c.Bind(&req); err != nil { //парсим JSON в структуру SumOfNumbers
		return c.JSON(http.StatusBadRequest, "Неверный формат JSON"+err.Error())
	}

	//Считаем сумму чисел
	var sum float64
	for _, num := range req.Numbers {
		sum += num
	}

	mu.Lock()
	totalSum += sum
	mu.Unlock()

	//Формируем ответ
	//response := map[string]float64{"sum": sum} //мапа для ответа
	return c.JSON(http.StatusOK, ResponseSumOfNumbers{ResponseNumbers: sum})

}

func GetSumHandler(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()
	return c.JSON(http.StatusOK, ResponseSumOfNumbers{ResponseNumbers: totalSum})
}
