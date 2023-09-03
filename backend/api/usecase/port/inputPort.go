package port

import "main/api/domain"

type InputPort interface {
	Input(domain.Koma) (domain.Koma, error)
}