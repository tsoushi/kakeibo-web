package resolver

import (
	"kakeibo-web-server/handler/graph/dataloader"
	"kakeibo-web-server/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	usecase *usecase.Usecase
	*dataloader.Loaders
}

func NewResolver(usecase *usecase.Usecase) *Resolver {
	return &Resolver{
		usecase: usecase,
		Loaders: dataloader.NewLoader(usecase),
	}
}
