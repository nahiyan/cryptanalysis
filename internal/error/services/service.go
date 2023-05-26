package services

import (
	"log"
)

func (errorSvc *ErrorService) Check(err error, message string) {
	if err != nil {
		log.Println(err)
		log.Println("Error:", message)
	}
}

func (errorSvc *ErrorService) Fatal(err error, message string) {
	if err != nil {
		log.Println(err)
		log.Fatal("Error:", message)
	}
}

func (errorSvc *ErrorService) Handle(err error, handler func(err error)) {
	if err != nil {
		handler(err)
	}
}
