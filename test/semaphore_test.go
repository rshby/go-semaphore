package test

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"testing"
	"time"
)

func TestSemaphoreChannel(t *testing.T) {
	var wg = &sync.WaitGroup{}
	var chanSemaphore = make(chan struct{}, 2)
	for i := 0; i < 10; i++ {
		chanSemaphore <- struct{}{}

		wg.Add(1)
		go func(wg *sync.WaitGroup, chanSemaphore chan struct{}) {
			defer func() {
				_ = <-chanSemaphore
				wg.Done()
			}()

			time.Sleep(1 * time.Second)
			fmt.Println("print from goroutine")
		}(wg, chanSemaphore)
	}

	wg.Wait()
	close(chanSemaphore)
}

func TestSemaphore(t *testing.T) {
	var (
		maximumAllowed = int64(2)
		semLimit       = semaphore.NewWeighted(maximumAllowed)
		wg             = &sync.WaitGroup{}
		ctx            = context.Background()
	)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			if err := semLimit.Acquire(ctx, 1); err != nil {
				continue
			}

			wg.Add(1)
			go func(wg *sync.WaitGroup, semLimit *semaphore.Weighted) {
				defer func() {
					semLimit.Release(1)
					wg.Done()
				}()

				time.Sleep(1 * time.Second)
				fmt.Println("print inside a goroutine")
			}(wg, semLimit)
		}
	}(wg)

	wg.Wait()
}
