package internal

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync/atomic"
	"syscall"

	"api-gateway/config"
	"api-gateway/internal/pkg/closer"

	"github.com/go-chi/chi/v5"
	"github.com/not-for-prod/clay/server"
	"github.com/not-for-prod/clay/transport"
	"google.golang.org/grpc"
)

// App application
type App struct {
	mainServer *server.Server

	mux *chi.Mux

	publicCloser *closer.Closer

	controllers []transport.ServiceDesc

	grpcConn map[string]grpc.ClientConnInterface

	started atomic.Bool
}

// New конструктор
func New(ctx context.Context) *App {
	app := &App{grpcConn: make(map[string]grpc.ClientConnInterface),
		publicCloser: closer.New(syscall.SIGTERM, syscall.SIGINT)}
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

	// start signal
	a.started.Store(true)

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
		a.initGrpcConn,
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
