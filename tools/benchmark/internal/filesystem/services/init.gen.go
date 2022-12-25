package services

import (
    "github.com/samber/do"
    
)

type FilesystemService struct {
    
}

func NewFilesystemService(i *do.Injector) (*FilesystemService, error) {

    svc := &FilesystemService{
    }

    

	return svc, nil
}
