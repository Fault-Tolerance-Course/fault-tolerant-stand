package internal

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"syscall"

	"order-service/config"

	"order-service/internal/application/service"
	"order-service/internal/infrastructure/dal"
	"order-service/internal/infrastructure/messagebus"
	"order-service/internal/pkg/closer"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/not-for-prod/clay/server"
	"github.com/not-for-prod/clay/transport"
)

// App application
type App struct {
	mainServer *server.Server

	mux *chi.Mux

	publicCloser *closer.Closer

	pool *pgxpool.Pool

	dal *dal.Registry

	services *service.Registry

	messageBus *messagebus.Registry

	controllers []transport.ServiceDesc
}

// New конструктор
func New(ctx context.Context) *App {
	app := &App{publicCloser: closer.New(syscall.SIGTERM, syscall.SIGINT)}
	err := app.init(ctx)
	if err != nil {
		log.Fatalf("[APP] Не удалось инициализировать приложение: %s", err.Error())
	}

	return app
}

// Run запуск приложения
func (a *App) Run(_ context.Context) {
	if a.mainServer != nil {
		go func() {
			if err := a.mainServer.Run(a.controllers...); err != nil {
				slog.Error(fmt.Sprintf("main server: %s", err.Error()))
				a.publicCloser.CloseAll()
			}
		}()
	}

	slog.Info(fmt.Sprintf("APP STARTED ON PORTS => GRPC: %d, HTTP: %d",
		config.Instance().GrpcServer.Port,
		config.Instance().HttpServer.Port,
	))

	a.publicCloser.Wait()

	closer.CloseAll()
}

func (a *App) init(ctx context.Context) error {
	initFuncs := []func(context.Context) error{
		a.initMainServer,
		a.initPostgres,
		a.initMessageBus,
		a.initDAL,
		a.initServices,
		a.initControllers,
	}

	for _, f := range initFuncs {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
