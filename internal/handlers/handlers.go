package handlers

import "sync"

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
