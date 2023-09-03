package port

import "main/api/domain"

type InputPort interface {
	GameMain(domain.Koma) (*domain.GameResponse, error)
}