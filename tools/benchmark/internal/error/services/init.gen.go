package services

import do "github.com/samber/do"

type ErrorService struct{}

func NewErrorService(injector *do.Injector) (*ErrorService, error) {
	svc := &ErrorService{}
	return svc, nil
}
