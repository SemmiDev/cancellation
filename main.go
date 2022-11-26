package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)  //  SIGINT is a ^C
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) //  SIGTERM is a kill

	ctx := context.Background()
	timeout := time.Second * 2
	ctx, stop := context.WithTimeout(ctx, timeout) // ctx will be canceled after 1 second
	defer stop()

	go func() {
		<-c
		fmt.Println("Received signal")
		fmt.Println("printNumbers will be canceled")
		stop()
		fmt.Println("canceled!")
	}()

	printNumbers(ctx, 10, func(i int) bool {
		return i%2 == 0
	})

	fmt.Println("Done")
}

type filter func(i int) bool

func printNumbers(ctx context.Context, n int, f filter) error {
	if isCtxCanceled(ctx) {
		return ctx.Err()
	}

	for i := 0; i < n; i++ {
		time.Sleep(time.Millisecond * 200)
		if isCtxCanceled(ctx) {
			return ctx.Err()
		}

		if f(i) {
			fmt.Println(i)
		}

		if isCtxCanceled(ctx) {
			return ctx.Err()
		}
	}
	return nil
}

func isCtxCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
