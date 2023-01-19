package services

import (
	"log"

	"github.com/sirupsen/logrus"
)

func (errorSvc *ErrorService) Check(err error, message string) {
	if err != nil {
		logrus.Error("Error:", message)
	}
}

func (errorSvc *ErrorService) Fatal(err error, message string) {
	if err != nil {
		logrus.Error(err)
		log.Fatal("Error:", message)
	}
}

func (errorSvc *ErrorService) Handle(err error, handler func(err error)) {
	if err != nil {
		handler(err)
	}
}
