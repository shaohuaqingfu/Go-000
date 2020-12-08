package main

import (
	"Week03/src/errgroup"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

func main() {
	fmt.Println("123")
	group := errgroup.WithCancel(context.Background()).WithPool(10)
	group.Go(func(ctx context.Context) error {
		return http.ListenAndServe("http://127.0.0.1:8080", nil)
	})
	group.Go(func(ctx context.Context) error {
		//return http.ListenAndServe("http://127.0.0.1:8081", nil)
		return errors.New("123")
	})
	fmt.Println("123")
	err := group.Wait()
	if err != nil {
		_ = fmt.Errorf("%+v", errors.Wrap(err, "启动服务出现异常"))
	}
}
