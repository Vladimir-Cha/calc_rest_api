package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	ResponseNumbers       float64 `json:"sum"`
	MultiplicationNumbers float64 `json:"multipl"`
}

var (
	totalSum float64
	totalMul float64
	mu       sync.Mutex // Переменная для mutex
)

// sumHandler - обработчик HTTP-запросов по пути /sum
func sumHandler(c echo.Context) error {
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

func getSumHandler(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()
	return c.JSON(http.StatusOK, ResponseSumOfNumbers{ResponseNumbers: totalSum})
}

func multimplHandler(c echo.Context) error {
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

func getMulHandler(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()
	return c.JSON(http.StatusOK, ResponseSumOfNumbers{MultiplicationNumbers: totalMul})
}

func main() {
	//Создаем сервер
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	//Регистрируем маршрут
	e.POST("/sum", sumHandler)
	e.GET("/totalsum", getSumHandler)
	e.POST("/mul", multimplHandler)
	e.GET("/totalmul", getMulHandler)

	fmt.Println("Запускаем сервер")

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal("Ошибка запуска сервера: ", err)
	}
}
