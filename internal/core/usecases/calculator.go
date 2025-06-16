package usecases

import (
	"github.com/Vladimir-Cha/calc_rest_api/internal/entities"
	"github.com/Vladimir-Cha/calc_rest_api/internal/interfaces"
	"github.com/google/uuid"
)

type Calculator struct {
	repo interfaces.MathRepository
}

func NewCalculator(repo interfaces.MathRepository) *Calculator {
	return &Calculator{repo: repo}
}

func (uc *Calculator) Result(token string, nums []float64) (entities.NumResponse, error) {
	//Считаем сумму и произведение чисел
	var sum float64
	var mul float64
	sum, mul = 0.0, 1.0

	for _, num := range nums {
		sum += num
		mul *= num
	}

	err := uc.repo.SaveResult(token, entities.NumResponse{
		Response_numbers:       sum,
		Multiplication_numbers: mul,
	})

	return entities.NumResponse{
		Response_numbers:       sum,
		Multiplication_numbers: mul,
	}, err
}

func (uc *Calculator) GetTotal() (entities.TotalResult, error) {
	return uc.repo.GetTotal()
}

func (uc *Calculator) GetResult(token string) (entities.NumResponse, entities.TotalResult, error) {
	res, err := uc.repo.GetResult(token)
	if err != nil {
		return entities.NumResponse{}, entities.TotalResult{}, err
	}

	total, err := uc.repo.GetTotal()
	return res, total, err
}

func (uc *Calculator) GenerateToken() string {
	return uuid.New().String()
}
