package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/Vladimir-Cha/calc_rest_api/docs"
	"github.com/Vladimir-Cha/calc_rest_api/internal/adapters/handlers"
	"github.com/Vladimir-Cha/calc_rest_api/internal/adapters/storage"
	"github.com/Vladimir-Cha/calc_rest_api/internal/core/usecases"
)

// @title Calculator REST API
// @version 1.0
// @description API для вычислений сумм и произведений

// @contact.name Vladimir
// @contact.email chaykovskyv@inbox.ru

// @host localhost:8080
// @BasePath /
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
	e.GET("/tokenresult", handlers.GetResult) //Get запрос по токену /totalres?token=XXX
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
