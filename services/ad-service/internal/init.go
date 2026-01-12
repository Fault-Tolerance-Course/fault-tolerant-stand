package internal

import (
	"context"
	"fmt"
	"log/slog"

	"ad-service/config"
	"ad-service/internal/app/ad/v1"
	"ad-service/internal/application/service"
	"ad-service/internal/infrastructure/gateway"
	"ad-service/internal/infrastructure/storage"
	"ad-service/internal/pkg/circuit"
	"ad-service/internal/pkg/closer"
	"ad-service/internal/pkg/connector/postgres"
	"ad-service/internal/pkg/grpc/intercept"
	"ad-service/internal/pkg/hedge"
	adV1 "ad-service/internal/pkg/pb/ad-service/ad/v1"

	"github.com/not-for-prod/clay/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/go-chi/chi/v5"
	"github.com/not-for-prod/clay/server"
)

func (a *App) initPostgres(ctx context.Context) error {
	pool, err := postgres.Pool(ctx, config.Instance().PostgresDSN())
	if err != nil {
		return fmt.Errorf("[POSTGRES] Не удалось инициализировать pool: %s", err.Error())
	}

	a.pool = pool
	return nil
}

func (a *App) initStorages(_ context.Context) error {
	if a.storages == nil {
		a.storages = storage.NewRegistry(a.pool)
	}
	return nil
}

func (a *App) initServices(_ context.Context) error {
	if a.services == nil {
		a.services = service.NewRegistry(a.storages, a.gateways)
	}
	return nil
}

func (a *App) initGateways(_ context.Context) error {
	if a.gateways == nil {
		a.gateways = gateway.NewRegistry(a.grpcConn)
	}
	return nil
}

func (a *App) initControllers(_ context.Context) error {
	a.controllers = []transport.ServiceDesc{
		adV1.NewAdServiceServiceDesc(ad.NewAdService(a.services)),
	}
	return nil
}

func (a *App) initMainServer(ctx context.Context) error {
	a.mux = chi.NewMux()
	// init server (http,grpc)
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
			grpc.ChainUnaryInterceptor(intercept.ErrorInterceptor()),
		),
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
			slog.Warn("review-service: main server gracefully stopped")
		case <-gracefulCtx.Done():
			err := fmt.Errorf("review-service: error while graceful shutdown server: %w", ctx.Err())
			_ = a.mainServer.Stop(ctx)
			return fmt.Errorf("review-service: stopped: %w", err)
		}
		return nil
	})

	return nil
}

func (a *App) initGrpcConn(_ context.Context) error {
	for _, srv := range []string{config.ReviewService} {
		var err error

		conn, err := grpc.NewClient(config.Instance().Targets[srv],
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				circuit.NewCircuitBreaker(config.Instance().Circuit).UnaryClientInterceptor(),
				hedge.NewHedger(config.Instance().Hedge).UnaryClientInterceptor(),
			),
		)

		if err != nil {
			return fmt.Errorf("не удалось инициализировать grpc соединение к %s : %s", srv, err.Error())
		}

		a.grpcConn[srv] = conn
		closer.Add(conn.Close)
	}
	return nil
}
