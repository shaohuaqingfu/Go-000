package main

import (
	"context"
	"net/http"
)

type A interface {
}

func main() {

	down := make(chan error)
	stop := make(chan struct{})

	go func() {
		down <- serverApp(stop)
	}()

	context.Background()

	for i := 0; i < cap(down); i++ {
		if err := <-down; err != nil {
			close(stop)
		}
	}
	for {
		select {
		case <-stop:
			break
		}
	}

}

func serverApp(stop chan struct{}) error {

	go func() {
		<-stop
		Shutdown()
	}()

	return http.ListenAndServe("127.0.0.1:8080", nil)
}

func Shutdown() {

}
