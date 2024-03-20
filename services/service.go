package services

import (
	"popwa/domain"
)
type popWaService struct {
	popWaRepo domain.Repositories

}

func NewPopWaService(popWaRepo domain.Repositories) *popWaService {
	return &popWaService{
		popWaRepo: popWaRepo,
	}
}