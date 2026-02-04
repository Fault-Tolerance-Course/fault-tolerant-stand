package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime"
	"sync/atomic"

	"api-gateway/config"
	"api-gateway/internal/pkg/grpc/intercept"
	"api-gateway/internal/pkg/healthcheck"

	"api-gateway/internal/pkg/http/middleware/timeout"

	"api-gateway/internal/pkg/retry"

	chimw "github.com/go-chi/chi/v5/middleware"

	"api-gateway/internal/app/ad/v1"
	"api-gateway/internal/app/order/v1"

	"api-gateway/internal/pkg/closer"

	adV1 "api-gateway/internal/pkg/pb/api-gateway/ad/v1"
	orderV1 "api-gateway/internal/pkg/pb/api-gateway/order/v1"

	externalAdV1 "api-gateway/internal/pkg/pb/external/ad-service/ad/v1"
	externalOrderV1 "api-gateway/internal/pkg/pb/external/order-service/order/v1"

	"github.com/go-chi/chi/v5"

	"github.com/not-for-prod/clay/server"
	"github.com/not-for-prod/clay/transport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func (a *App) initControllers(_ context.Context) error {
	a.controllers = []transport.ServiceDesc{
		adV1.NewAdServiceServiceDesc(ad.NewAdService(
			externalAdV1.NewAdServiceClient(a.grpcConn[config.AdService]),
		)),
		orderV1.NewOrderServiceServiceDesc(order.NewOrderService(
			externalOrderV1.NewOrderServiceClient(a.grpcConn[config.OrderService]),
		)),
	}
	return nil
}

func (a *App) initMainServer(ctx context.Context) error {
	a.mainMux = chi.NewMux()

	// используем middleware для установки timeout'а на обработку запроса
	a.mainMux.Use(timeout.Middleware)

	// init server (htt,grpc)
	a.mainServer = server.NewServer(
		config.Instance().GrpcServer.Port,
		server.WithHTTPMux(a.mainMux),
		server.WithHTTPPort(config.Instance().HttpServer.Port),
		server.WithGRPCOpts(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: config.Instance().GrpcServer.MaxConnectionIdle,
				MaxConnectionAge:  config.Instance().GrpcServer.MaxConnectionAge,
				Time:              config.Instance().GrpcServer.Time,
				Timeout:           config.Instance().GrpcServer.Timeout,
			}),
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
			slog.Warn("api-gateway: main server gracefully stopped")
		case <-gracefulCtx.Done():
			err := fmt.Errorf("api-gateway: error while graceful shutdown server: %w", gracefulCtx.Err())
			_ = a.mainServer.Stop(ctx) // TODO: поправить в либе на hard shutdown (да, заметил поздно :) )
			return fmt.Errorf("api-gateway: stopped: %w", err)
		}
		return nil
	})

	return nil
}

func (a *App) initAdminServer(ctx context.Context) error {
	// init admin listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Instance().HttpServer.AdminPort))
	if err != nil {
		return fmt.Errorf("failed to init admin listener: %w", err)
	}

	a.adminListener = lis

	// init healthcheck for admin server
	if err = a.initHealthCheck(ctx); err != nil {
		return fmt.Errorf("failed to init healthcheck: %w", err)
	}

	a.adminMux = chi.NewMux()

	a.adminMux.Mount("/debug", chimw.Profiler())

	// register healthcheck
	a.adminMux.HandleFunc(healthcheck.LivenessPath, a.healthCheck.LiveEndpoint)
	a.adminMux.HandleFunc(healthcheck.ReadinessPath, a.healthCheck.ReadyEndpoint)

	return nil
}

func (a *App) initGrpcConn(_ context.Context) error {
	for _, srv := range []string{config.AdService, config.OrderService} {
		var err error

		conn, err := grpc.NewClient(config.Instance().Targets[srv],
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				intercept.SetClientNameInterceptor(config.AppName),
				retry.NewRetry(config.Instance().Retry).UnaryClientInterceptor()),
		)

		if err != nil {
			return fmt.Errorf("не удалось инициализировать grpc соединение к %s : %s", srv, err.Error())
		}

		a.grpcConn[srv] = conn
		closer.Add(conn.Close)
	}
	return nil
}

func (a *App) initHealthCheck(_ context.Context) error {
	a.healthCheck = healthcheck.NewHandler()

	// поверяю, что нет утечки горутин на старте (как пример)
	a.healthCheck.AddLivenessCheck("goroutines", func() error {
		if runtime.NumGoroutine() < 1000 {
			return nil
		}
		return fmt.Errorf("application has too much running goroutines")
	})

	// readiness - т.к. я уже проинициализировал все компоненты
	a.healthCheck.AddReadinessCheck("started", func() error {
		if atomic.LoadInt32(&a.started) != 0 {
			return nil
		}
		return fmt.Errorf("application is not statred yet")
	})

	a.healthCheck.AddReadinessCheck("termination", func() error {
		if atomic.LoadInt32(&a.terminated) == 0 {
			return nil
		}
		return fmt.Errorf("application is terminating now")
	})

	a.publicCloser.Add(func() error {
		slog.Warn(fmt.Sprintf("app got termination signal, graceful config timeout: %s",
			config.Instance().Graceful.Timeout.String()))

		atomic.StoreInt32(&a.terminated, 1)
		return nil
	})
	return nil
}
