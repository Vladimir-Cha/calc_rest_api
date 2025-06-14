package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Vladimir-Cha/calc_rest_api/internal/adapters/handlers"
	"github.com/Vladimir-Cha/calc_rest_api/internal/adapters/storage"
	"github.com/Vladimir-Cha/calc_rest_api/internal/core/usecases"
)

func main() {
	//Инициируем хранилище
	storage := storage.NewMathStorage()
	//Use cases
	calculator := usecases.NewCalculator(storage)
	// HTTP-обработчики
	handlers := handlers.NewHandlers(calculator)
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
