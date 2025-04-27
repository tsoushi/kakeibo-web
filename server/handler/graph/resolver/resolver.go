package resolver

import "kakeibo-web-server/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	usecase *usecase.Usecase
}

func NewResolver(usecase *usecase.Usecase) *Resolver {
	return &Resolver{
		usecase: usecase,
	}
}
