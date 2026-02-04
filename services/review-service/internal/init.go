package internal

import (
	"context"
	"fmt"
	"log/slog"

	"review/config"
	review "review/internal/app/review/v1"
	"review/internal/pkg/grpc/intercept"
	"review/internal/pkg/loadshedding"
	reviewV1 "review/internal/pkg/pb/review-service/review/v1"

	"review/internal/application/service"
	"review/internal/infrastructure/storage"

	"github.com/go-chi/chi/v5"
	"github.com/not-for-prod/clay/server"
	"github.com/not-for-prod/clay/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func (a *App) initStorages(_ context.Context) error {
	if a.storages == nil {
		a.storages = storage.NewRegistry()
	}
	return nil
}

func (a *App) initServices(_ context.Context) error {
	if a.services == nil {
		a.services = service.NewRegistry(a.storages)
	}
	return nil
}

func (a *App) initControllers(_ context.Context) error {
	a.controllers = []transport.ServiceDesc{
		reviewV1.NewReviewServiceServiceDesc(review.NewReviewService(a.services)),
	}
	return nil
}

func (a *App) initMainServer(ctx context.Context) error {
	a.mux = chi.NewMux()
	// init server (htt,grpc)
	a.mainServer = server.NewServer(
		config.Instance().GrpcServer.Port,
		server.WithHTTPMux(a.mux),
		server.WithHTTPPort(config.Instance().HttpServer.Port),
		server.WithGRPCOpts(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: config.Instance().GrpcServer.MaxConnectionIdle,
				MaxConnectionAge:  config.Instance().GrpcServer.MaxConnectionAge,
				Time:              config.Instance().GrpcServer.Time,
				Timeout:           config.Instance().GrpcServer.Timeout,
			}),
			grpc.ChainUnaryInterceptor(
				intercept.ErrorInterceptor(),
				//intercept.HedgedDemoInterceptor(),
				intercept.CircuitDemoInterceptor(),
				loadshedding.NewLoadShedding(config.Instance().LoadShedding).UnaryServerInterceptor(),
			),
		),
	)

	a.publicCloser.Add(func() error {
		gracefulCtx, cancel := context.WithTimeout(context.Background(), config.Instance().Graceful.Timeout)
		defer cancel()

		done := make(chan struct{})
		go func() {
			err := a.mainServer.Stop(gracefulCtx)
			if err != nil {
				slog.Error(fmt.Sprintf("stop main server error: %s", err.Error()))
			}
			close(done)
		}()

		select {
		case <-done:
			slog.Warn("review-service: main server gracefully stopped")
		case <-gracefulCtx.Done():
			err := fmt.Errorf("review-service: error while graceful shutdown server: %w", gracefulCtx.Err())
			_ = a.mainServer.Stop(ctx) // TODO: поправить в либе на hard shutdown (да, заметил поздно :) )
			return fmt.Errorf("review-service: stopped: %w", err)
		}
		return nil
	})

	return nil
}
