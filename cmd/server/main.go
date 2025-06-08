package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/Vladimir-Cha/calc_rest_api/internal/handlers"
)

func main() {
	//Создаем сервер
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	//Регистрируем маршрут
	e.POST("/result", handlers.Result)
	e.GET("/totalresult", handlers.GetTotal)
	e.GET("/totalres", handlers.GetResult) //Get запрос по токену /totalres?token=XXX

	e.POST("/token", handlers.Token)

	log.Println("Запускаем сервер")

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":8080"
	}

	if err := e.Start(port); err != nil {
		e.Logger.Fatal("Ошибка запуска сервера: ", err)
	}
}
