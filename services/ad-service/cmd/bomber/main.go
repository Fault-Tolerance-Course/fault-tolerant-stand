package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	adV1 "ad-service/internal/pkg/pb/ad-service/ad/v1"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:7002", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := adV1.NewAdServiceClient(conn)

	// Создаём объявление 1 раз
	resp, err := client.CreateAd(context.Background(), &adV1.CreateAdRequest{
		Title:    "title",
		Category: "category",
		AuthorId: uuid.New().String(),
		Price:    &money.Money{Units: 100},
	})
	if err != nil {
		panic(err)
	}

	adID := resp.GetAd().GetId()

	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"x-client-name", "api-gateway",
	)

	// Тестовые параметры
	phases := []struct {
		rps      int
		duration time.Duration
	}{
		{rps: 50, duration: 3 * time.Second}, // > limit, чтобы показать исчерпание
		{rps: 0, duration: 2 * time.Second},  // пауза, токены копятся
		{rps: 15, duration: 5 * time.Second}, // < limit, все проходят
	}

	for i, p := range phases {
		fmt.Printf("\n--- phase %d: rps=%d duration=%s ---\n", i+1, p.rps, p.duration)

		runPhase(ctx, client, adID, p.rps, p.duration)
	}
}

func runPhase(ctx context.Context, client adV1.AdServiceClient, adID string, rps int, duration time.Duration) {
	if rps == 0 {
		time.Sleep(duration)
		return
	}

	var (
		wg      sync.WaitGroup
		ticker  = time.NewTicker(time.Second / time.Duration(rps))
		timeout = time.After(duration)
	)

	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			wg.Wait()
			return

		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := client.GetAd(ctx, &adV1.GetAdRequest{AdId: adID})
				if err != nil {
					fmt.Println("[FAIL]", err)
				} else {
					fmt.Println("[OK]")
				}
			}()
		}
	}
}
