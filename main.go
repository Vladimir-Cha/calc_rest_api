package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*Калькулятор работает в Postman. Нужно сделать Post запрос по адресу http://localhost:8080/sum
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

const contentTypeJSON = "application/JSON"

// sumHandler - обработчик HTTP-запросов по пути /sum
func sumHandler(w http.ResponseWriter, r *http.Request) {
	//Если метод не POST, то выводим ошибку
	if r.Method != http.MethodPost {
		http.Error(w, "Метод должен быть POST", http.StatusMethodNotAllowed)
		return
	}

	var req SumOfNumbers
	decoder := json.NewDecoder(r.Body)           //декодер для тела запроса
	if err := decoder.Decode(&req); err != nil { //парсим JSON в структуру SumOfNumbers
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest) //если ошибка парсинга
		return
	}

	//Считаем сумму чисел
	var sum float64
	for _, num := range req.Numbers {
		sum += num
	}

	//Формируем ответ
	response := map[string]float64{"sum": sum} //мапа для ответа
	w.Header().Set("Content-type", contentTypeJSON)

	json.NewEncoder(w).Encode(response)

}

func main() {
	http.HandleFunc("/sum", sumHandler)
	fmt.Println("Запускаем сервер")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
