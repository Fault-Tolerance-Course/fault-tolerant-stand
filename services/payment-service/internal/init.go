package internal

import (
	"context"
	"fmt"
	"log/slog"

	"payment-service/config"
	"payment-service/internal/application/service"
	"payment-service/internal/infrastructure/dal"
	"payment-service/internal/infrastructure/messagebus"
	"payment-service/internal/pkg/connector/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/not-for-prod/clay/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func (a *App) initControllers(_ context.Context) error {
	// init service
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

func (a *App) initDAL(_ context.Context) error {
	if a.dal == nil {
		a.dal = dal.NewRegistry(a.pool)
	}
	return nil
}

func (a *App) initServices(_ context.Context) error {
	if a.services == nil {
		a.services = service.NewRegistry(a.dal)
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
			})),
	)

	a.publicCloser.Add(func() error {
		gracefulCtx, cancel := context.WithTimeout(context.Background(), config.Instance().Graceful.Timeout)
		defer cancel()

		done := make(chan struct{})
		go func() {
			err := a.mainServer.Stop(ctx)
			if err != nil {
				slog.Error(fmt.Sprintf("stop main server error: %s", err.Error()))
			}
			close(done)
		}()

		select {
		case <-done:
			slog.Warn("order-service: main server gracefully stopped")
		case <-gracefulCtx.Done():
			err := fmt.Errorf("order-service: error while graceful shutdown server: %w", ctx.Err())
			_ = a.mainServer.Stop(ctx)
			return fmt.Errorf("order-service: stopped: %w", err)
		}
		return nil
	})

	return nil
}
func (a *App) initMessageBus(_ context.Context) error {
	if a.messageBus == nil {
		a.messageBus = messagebus.NewRegistry(a.services)
	}
	return nil
}
