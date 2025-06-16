package storage

import (
	"errors"
	"sync"

	"github.com/Vladimir-Cha/calc_rest_api/internal/entities"
)

// структура для хранения результатов
type MathStorage struct {
	mu      sync.RWMutex
	results map[string]entities.NumResponse
	total   entities.TotalResult
}

func NewMathStorage() *MathStorage {
	return &MathStorage{
		results: make(map[string]entities.NumResponse),
	}
}

func (s *MathStorage) SaveResult(token string, res entities.NumResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[token] = res
	s.total.TotalSum += res.ResponseNumbers
	s.total.TotalMul += res.MultiplicationNumbers
	return nil
}

func (s *MathStorage) GetResult(token string) (entities.NumResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, ok := s.results[token]
	if !ok {
		return entities.NumResponse{}, errors.New("Токен не найден")
	}
	return res, nil
}

func (s *MathStorage) GetTotal() (entities.TotalResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.total, nil
}
