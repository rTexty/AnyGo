package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

const WORKERS int = 20

func main(){
	var wg sync.WaitGroup
	
	jobs := make(chan int, 100)
	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for range WORKERS{
		wg.Add(1)
	go worker(jobs, &wg)

	}

	go func(){
		defer close(jobs)
		for i := range 100{
			select{
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- i:
			}
		}
	}()

	wg.Wait()

	
}

func worker(ch <- chan int, wg *sync.WaitGroup){
	defer wg.Done()
	for item := range ch{
		_ = item

	}
}
