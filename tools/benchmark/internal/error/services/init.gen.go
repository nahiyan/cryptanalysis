package services

import (
    "github.com/samber/do"
    
)

type ErrorService struct {
    
}

func NewErrorService(i *do.Injector) (*ErrorService, error) {

    svc := &ErrorService{
    }

    

	return svc, nil
}
