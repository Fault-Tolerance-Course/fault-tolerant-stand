package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"payment-service/config"
	"payment-service/internal/application/service"
	"payment-service/internal/infrastructure/dal"
	"payment-service/internal/pkg/closer"

	order_created "payment-service/internal/infrastructure/inbox/order-events/order-created"

	orderevents "payment-service/internal/infrastructure/inbox/order-events"
	"payment-service/internal/infrastructure/messagebus"
	"payment-service/internal/pkg/connector/postgres"
	event_router "payment-service/internal/pkg/event-router"
	"payment-service/internal/pkg/inbox"
	"payment-service/internal/pkg/leader"
	"payment-service/internal/pkg/worker"

	"github.com/go-chi/chi/v5"
	"github.com/not-for-prod/clay/server"
	clientV3 "go.etcd.io/etcd/client/v3"
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

func (a *App) initEtcdClient(ctx context.Context) error {
	client, err := clientV3.New(clientV3.Config{
		Endpoints:   config.Instance().ETCD.Endpoints,
		DialTimeout: config.Instance().ETCD.DialTimeout,
		Context:     ctx,
	})
	if err != nil {
		return fmt.Errorf("[ETCD] Не удалось инициализировать кликента: %s", err.Error())
	}

	a.etcdClient = client
	return nil
}

func (a *App) initDAL(_ context.Context) error {
	if a.dal == nil {
		a.dal = dal.NewRegistry(a.pool)
	}
	return nil
}

func (a *App) initInbox(_ context.Context) error {
	a.inbox = inbox.NewInbox(a.pool, config.Instance().Inbox.Config, inbox.WithMetrics())

	return nil
}

func (a *App) initHandlers(_ context.Context) error {
	a.eventRouter = event_router.NewEventRouter[string, []byte]()

	a.eventRouter.RegisterAll(order_created.NewHandler(a.services.Payment))

	a.inbox.RegisterHandler(orderevents.NewHandler(
		config.OrderEventsTopic,
		config.Instance().InboxConfig(config.OrderEventsTopic).BatchSize,
		a.eventRouter,
	))
	return nil
}

func (a *App) initElectionManager(ctx context.Context) error {
	a.electionManager = leader.NewElectionManager(a.etcdClient, config.Instance().LeaderElection.Key)

	a.electionManager.AddWorker(worker.NewWorker(ctx,
		a.inbox.HandlePendingMessages(inbox.ModeNormal),
		func(ctx context.Context) time.Duration {
			return config.Instance().Inbox.Config.NormalMessagesPollInterval
		}, func(ctx context.Context) int { return 1 }),
	)

	a.electionManager.AddWorker(worker.NewWorker(ctx,
		a.inbox.HandlePendingMessages(inbox.ModeError),
		func(ctx context.Context) time.Duration {
			return config.Instance().Inbox.Config.ErrorMessagesPollInterval
		}, func(ctx context.Context) int { return 1 }),
	)

	closer.Add(func() error {
		a.electionManager.Stop()
		return nil
	})
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
			err := a.mainServer.Stop(gracefulCtx)
			if err != nil {
				slog.Error(fmt.Sprintf("stop main server error: %s", err.Error()))
			}
			close(done)
		}()

		select {
		case <-done:
			slog.Warn("payment-service: main server gracefully stopped")
		case <-gracefulCtx.Done():
			err := fmt.Errorf("payment-service: error while graceful shutdown server: %w", gracefulCtx.Err())
			_ = a.mainServer.Stop(ctx) // TODO: поправить в либе на hard shutdown (да, заметил поздно :) )
			return fmt.Errorf("payment-service: stopped: %w", err)
		}
		return nil
	})

	return nil
}
func (a *App) initMessageBus(_ context.Context) error {
	if a.messageBus == nil {
		a.messageBus = messagebus.NewRegistry(a.inbox)
	}
	return nil
}
