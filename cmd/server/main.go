package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	serverErr := make(chan error, 1)
	go func() {
		port := os.Getenv("TODO_PORT")
		if port == "" {
			port = ":8080"
		}
		log.Println("Запускаем сервер")
		if err := e.Start(port); err != nil {
			serverErr <- err
		}
	}()

	// канал для gracefull shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Printf("Ошибка сервера: %v", err)
	case <-shutdown:
		log.Println("Shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal("Ошибка при завершении сервера: ", err)
		}
		log.Println("Shutdown complete")
	}

}
