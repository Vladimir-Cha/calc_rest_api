package entities

// Структура SumOfNumbers для входящего JSON-запроса
type NumRequest struct {
	Numbers []float64 `json:"numbers"`
}

type NumResponse struct {
	ResponseNumbers       float64 `json:"sum"`
	MultiplicationNumbers float64 `json:"multipl"`
}

type TotalResult struct {
	TotalSum float64 `json:"totalsum"`
	TotalMul float64 `json:"totalmultipl"`
}
