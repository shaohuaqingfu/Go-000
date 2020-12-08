package sync

import (
	"context"
	"log"
)

func Go(ctx context.Context, f func(ctx context.Context)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[系统错误] error = %+v", err)
			}
		}()
		f(ctx)
		select {
		case <-ctx.Done():
			log.Printf("[系统错误] error = %+v\n", ctx.Err())
			return
		}
	}()
}
