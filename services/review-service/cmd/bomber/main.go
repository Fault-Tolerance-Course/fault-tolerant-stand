package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	reviewV1 "review/internal/pkg/pb/review-service/review/v1"

	"google.golang.org/grpc"
)

func main() {
	decreaseLimit()
	//increaseLimit()
}

func decreaseLimit() {
	// Подключение к gRPC серверу
	conn, err := grpc.Dial("localhost:7003", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := reviewV1.NewReviewServiceClient(conn)

	var wg sync.WaitGroup
	totalRequests := 5000 // общее количество запросов
	concurrency := 200    // сколько одновременно горутин

	sem := make(chan struct{}, concurrency) // ограничение параллельных запросов

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{} // блокируем, если достигли concurrency

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }() // освобождаем слот

			// Контекст с таймаутом меньше, чем максимальная задержка метода
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			start := time.Now()
			_, err := client.UnstableMethod(ctx, &reviewV1.UnstableMethodRequest{})
			elapsed := time.Since(start)

			if err != nil {
				fmt.Printf("[Request %d] ERROR: %v (elapsed: %v)\n", i, err, elapsed)
			} else {
				fmt.Printf("[Request %d] OK (elapsed: %v)\n", i, elapsed)
			}
		}(i)

		// небольшой интервал между стартами, чтобы нагрузка была плавной
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("Нагрузка завершена")
}

func increaseLimit() {
	conn, _ := grpc.Dial("localhost:7003", grpc.WithInsecure())
	defer func() {
		_ = conn.Close()
	}()

	client := reviewV1.NewReviewServiceClient(conn)

	var inflight int64 = 1
	var successes int64
	var failures int64

	for {
		current := atomic.LoadInt64(&inflight)

		for i := 0; i < int(current); i++ {
			go func() {
				start := time.Now()
				_, err := client.UnstableMethod(context.Background(), &reviewV1.UnstableMethodRequest{})
				elapsed := time.Since(start)

				if err != nil {
					atomic.AddInt64(&failures, 1)
					fmt.Println(err.Error())
				} else {
					atomic.AddInt64(&successes, 1)

					// если быстро — увеличиваем давление
					if elapsed < 40*time.Millisecond {
						atomic.AddInt64(&inflight, 1)
					}
				}
			}()
		}

		// если слишком много ошибок — сбавим
		if atomic.LoadInt64(&failures) > 10 && atomic.LoadInt64(&inflight) > 1 {
			atomic.AddInt64(&inflight, -1)
			atomic.StoreInt64(&failures, 0)
		}

		fmt.Printf("inflight=%d ok=%d fail=%d\n",
			atomic.LoadInt64(&inflight),
			atomic.LoadInt64(&successes),
			atomic.LoadInt64(&failures),
		)

		time.Sleep(100 * time.Millisecond)
	}
}
