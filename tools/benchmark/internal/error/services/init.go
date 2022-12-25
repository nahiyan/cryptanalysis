package services

import "github.com/samber/do"

func NewErrorService(i *do.Injector) (*ErrorService, error) {
	return &ErrorService{}, nil
}
