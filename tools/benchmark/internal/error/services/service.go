package services

import (
	"fmt"
	"log"
)

func (errorSvc *ErrorService) Check(err error, handler func(error) string) {
	if err != nil {
		fmt.Println("Error:", handler(err))
	}
}

func (errorSvc *ErrorService) CheckWithFatal(err error, message string) {
	if err != nil {
		log.Fatal("Error:", message)
	}
}

func (errorSvc *ErrorService) CheckWithPanic(err error, handler func(error) string) {
	if err != nil {
		panic("Error:" + handler(err))
	}
}
