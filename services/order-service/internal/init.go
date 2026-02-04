package internal

import (
	"context"
	"fmt"
	"log/slog"

	"order-service/config"
	v1 "order-service/internal/app/order/v1"
	"order-service/internal/application/service"
	"order-service/internal/infrastructure/dal"
	"order-service/internal/infrastructure/messagebus"
	"order-service/internal/pkg/connector/postgres"
	"order-service/internal/pkg/grpc/intercept"
	orderV1 "order-service/internal/pkg/pb/order-service/order/v1"

	"github.com/go-chi/chi/v5"

	"github.com/not-for-prod/clay/server"
	"github.com/not-for-prod/clay/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func (a *App) initControllers(_ context.Context) error {
	a.controllers = []transport.ServiceDesc{
		orderV1.NewOrderServiceServiceDesc(v1.NewOrderService(a.services)),
	}
	return nil
}

func (a *App) initPostgres(ctx context.Context) error {
	pool, err := postgres.Pool(ctx, config.Instance().PostgresDSN())
	if err != nil {
		return fmt.Errorf("[POSTGRES] Не удалось инициализировать pool: %s", err.Error())
	}

	a.pool = pool
	return nil
}

func (a *App) initMessageBus(_ context.Context) error {
	if a.messageBus == nil {
		a.messageBus = messagebus.NewRegistry()
	}
	return nil
}

func (a *App) initDAL(_ context.Context) error {
	if a.dal == nil {
		a.dal = dal.NewRegistry(a.pool)
	}
	return nil
}

func (a *App) initServices(_ context.Context) error {
	if a.services == nil {
		a.services = service.NewRegistry(a.dal, a.messageBus)
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
				intercept.ExtractClientNameInterceptor(),
				intercept.ErrorInterceptor(),
				intercept.RetryDemoInterceptor(), // эмуляция сбоя для retry
			)),
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
			slog.Warn("order-service: main server gracefully stopped")
		case <-gracefulCtx.Done():
			err := fmt.Errorf("order-service: error while graceful shutdown server: %w", gracefulCtx.Err())
			_ = a.mainServer.Stop(ctx) // TODO: поправить в либе на hard shutdown (да, заметил поздно :) )
			return fmt.Errorf("order-service: stopped: %w", err)
		}
		return nil
	})

	return nil
}
