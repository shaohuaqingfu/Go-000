package main

import (
	"Week02/src/api"
	"log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover error: %s", err)
		}
	}()
	err := api.Init()
	if err != nil {
		log.Printf("error: %+v", err)
		return
	}
}

type errString struct {
	s string
}

func New(msg string) errString {
	return errString{s: msg}
}

func (e *errString) Error() string {
	return e.s
}
