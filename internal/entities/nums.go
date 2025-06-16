package entities

// Структура SumOfNumbers для входящего JSON-запроса
type NumRequest struct {
	Numbers []float64 `json:"numbers"`
}

type NumResponse struct {
	Response_numbers       float64 `json:"sum"`
	Multiplication_numbers float64 `json:"multipl"`
}

type TotalResult struct {
	Total_sum float64 `json:"totalsum"`
	Total_mul float64 `json:"totalmultipl"`
}
