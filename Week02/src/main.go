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
