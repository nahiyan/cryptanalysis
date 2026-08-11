package services

import do "github.com/samber/do"

type ConfigService struct {
	Properties
}

func NewConfigService(injector *do.Injector) (*ConfigService, error) {
	svc := &ConfigService{}
	svc.Init()
	return svc, nil
}
