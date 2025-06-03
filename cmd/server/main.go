package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/Vladimir-Cha/calc_rest_api/internal/handlers"
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

func main() {
	//Создаем сервер
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	//Регистрируем маршрут
	e.POST("/sum", handlers.SumHandler)
	e.GET("/totalsum", handlers.GetSumHandler)
	e.POST("/mul", handlers.MultimplHandler)
	e.GET("/totalmul", handlers.GetMulHandler)

	fmt.Println("Запускаем сервер")

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal("Ошибка запуска сервера: ", err)
	}
}
