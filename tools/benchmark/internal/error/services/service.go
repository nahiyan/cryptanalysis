package services

import (
	"fmt"
	"log"
)

func (errorSvc *ErrorService) Check(err error, message string) {
	if err != nil {
		fmt.Println("Error:", message)
	}
}

func (errorSvc *ErrorService) Fatal(err error, message string) {
	if err != nil {
		log.Fatal("Error:", message)
	}
}

// func (errorSvc *ErrorService) Panic(err error, handler func(error) string) {
// 	if err != nil {
// 		panic("Error:" + handler(err))
// 	}
// }
