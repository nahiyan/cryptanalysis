package services

import (
    "github.com/samber/do"
    
)

type ConfigService struct {
    Properties
}

func NewConfigService(i *do.Injector) (*ConfigService, error) {

    svc := &ConfigService{
    }

    svc.Init()

	return svc, nil
}
