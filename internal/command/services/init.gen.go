package services

import do "github.com/samber/do"

type CommandService struct{}

func NewCommandService(injector *do.Injector) (*CommandService, error) {
	svc := &CommandService{}
	return svc, nil
}
