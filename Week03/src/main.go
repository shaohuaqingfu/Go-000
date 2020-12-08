package main

import (
	"Week03/src/errgroup"
	"Week03/src/sync"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// 1.基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。
func main() {

	//exitChan := make(chan bool, 1)
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()
	group := errgroup.WithCancel(ctx).WithPool(10)

	server1 := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	server2 := &http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: nil,
	}
	group.Go(func(ctx context.Context, cancel context.CancelFunc) error {
		for s := range signalChan {
			switch s {
			case syscall.SIGINT:
				fmt.Printf("receive ctrl^c")
				cancel()
				break
			case syscall.SIGTERM:
				fmt.Printf("receive ctrl^\\")
				cancel()
				break
			case syscall.SIGQUIT:
				fmt.Printf("receive ctrl^\\")
				cancel()
				break
			}
		}
		return nil
	})
	group.Go(func(ctx context.Context, cancel context.CancelFunc) error {
		sync.Go(ctx, func(ctx context.Context) {
			var err error
			select {
			case <-ctx.Done():
				err = server1.Shutdown(ctx)
				//case <-exitChan:
				//	err = server1.Shutdown(ctx)
			}
			if err != nil {
				fmt.Printf("[系统错误] error = %+v\n", errors.Wrap(err, "server1 shutdown failed"))
			}
		})
		return server1.ListenAndServe()
	})
	group.Go(func(ctx context.Context, cancel context.CancelFunc) error {
		sync.Go(ctx, func(ctx context.Context) {
			var err error
			select {
			case <-ctx.Done():
				err = server2.Shutdown(ctx)
				//case <-exitChan:
				//	err = server2.Shutdown(ctx)
			}
			if err != nil {
				fmt.Printf("[系统错误] error = %+v\n", errors.Wrap(err, "server2 shutdown failed"))
			}
		})
		return server2.ListenAndServe()
	})
	err := group.Wait()
	if err != nil {
		fmt.Printf("[系统错误] error = %+v\n", errors.Wrap(err, "server error"))
	}
}
