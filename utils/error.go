package utils

import "log"

func HandleErr(message string, err error) {
	if err != nil {
		log.Fatal(message, err.Error())
	}
}

func PanicError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
