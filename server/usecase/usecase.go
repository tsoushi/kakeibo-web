package usecase

import "kakeibo-web-server/repository"

type Usecase struct {
	repo *repository.Repository
}

func NewUsecase(repository *repository.Repository) *Usecase {
	return &Usecase{
		repo: repository,
	}
}
