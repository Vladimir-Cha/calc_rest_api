package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

/*Калькулятор работает в любом web-клиенте.
Пример использования в Postman:
Нужно сделать Post запрос по адресу http://localhost:8080/sum
В Headers добавить ключ Content-Type, значение application/json
Тело запроса должно выглядеть следующим образом:

{
"numbers": [5, 10.2, 63]
}

5, 10.2, 63 - примеры чисел, сумму которых хотим получить

Ответ JSON будет выглядеть так:

{
    "sum": 78.2
}

78,2 - пример ответа суммы чисел из запроса
*/

// Структура SumOfNumbers для входящего JSON-запроса
type SumOfNumbers struct {
	Numbers []float64 `json:"numbers"`
}

type ResponseSumOfNumbers struct {
	ResponseNumbers float64 `json:"sum"`
}

// sumHandler - обработчик HTTP-запросов по пути /sum
func sumHandler(c echo.Context) error {
	//Если метод не POST, то выводим ошибку
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, "Метод должен быть POST")
	}

	var req SumOfNumbers
	//декодер для тела запроса
	if err := c.Bind(&req); err != nil { //парсим JSON в структуру SumOfNumbers
		return c.JSON(http.StatusBadRequest, "Неверный формат JSON")
	}

	//Считаем сумму чисел
	var sum float64
	for _, num := range req.Numbers {
		sum += num
	}

	//Формируем ответ
	//response := map[string]float64{"sum": sum} //мапа для ответа
	return c.JSON(http.StatusOK, ResponseSumOfNumbers{ResponseNumbers: sum})

}

func main() {
	//Создаем сервер
	e := echo.New()

	//Регистрируем маршрут
	e.POST("/sum", sumHandler)

	fmt.Println("Запускаем сервер")

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal("Ошибка запуска сервера: ", err)
	}
}
