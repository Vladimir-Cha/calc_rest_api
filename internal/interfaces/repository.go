package interfaces

import "github.com/Vladimir-Cha/calc_rest_api/internal/entities"

type MathRepository interface {
	SaveResult(token string, res entities.NumResponse) error
	GetResult(token string) (entities.NumResponse, error)
	GetTotal() (entities.TotalResult, error)
}
