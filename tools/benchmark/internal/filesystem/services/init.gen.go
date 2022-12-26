package services

import do "github.com/samber/do"

type FilesystemService struct{}

func NewFilesystemService(injector *do.Injector) (*FilesystemService, error) {
	svc := &FilesystemService{}
	return svc, nil
}
