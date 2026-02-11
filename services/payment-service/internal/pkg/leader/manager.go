package leader

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	clientV3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Worker interface {
	Start(ctx context.Context) error
	Stop() error
}

type ElectionManager interface {
	AddWorker(w Worker)
	Run(ctx context.Context)
	Stop()
}

type manager struct {
	etcdClient *clientV3.Client

	electionKey string

	workers []Worker

	ctx    context.Context
	cancel func()

	done chan struct{}
}

func NewElectionManager(cli *clientV3.Client, electionKey string) ElectionManager {
	return &manager{
		etcdClient:  cli,
		electionKey: electionKey,
		done:        make(chan struct{}),
	}
}

func (m *manager) Run(ctx context.Context) {
	m.ctx, m.cancel = context.WithCancel(ctx)
	go m.loop()
}

func (m *manager) loop() {
	defer close(m.done)

	for {
		select {
		case <-m.ctx.Done():
			return
		default:
		}

		if err := m.electionCycle(); err != nil {
			slog.Error("election cycle failed", "err", err)
			time.Sleep(2 * time.Second)
		}
	}
}

func (m *manager) electionCycle() error {
	session, err := concurrency.NewSession(
		m.etcdClient,
		concurrency.WithTTL(15),
		concurrency.WithContext(m.ctx),
	)
	if err != nil {
		return err
	}
	defer session.Close()

	e := concurrency.NewElection(session, m.electionKey)

	slog.Info("campaigning for leadership")

	if err = e.Campaign(m.ctx, instanceID()); err != nil {
		return err
	}

	slog.Info("leadership acquired")

	leaderCtx, cancel := context.WithCancel(m.ctx)

	// следим за смертью сессии
	go func() {
		<-session.Done()
		slog.Warn("leadership lost (session expired)")
		cancel()
	}()

	m.onLeadingStart(leaderCtx)

	<-leaderCtx.Done()

	go m.onLeadingStop()

	return nil
}

func (m *manager) AddWorker(w Worker) {
	m.workers = append(m.workers, w)
}

func (m *manager) onLeadingStart(ctx context.Context) {
	slog.Info("[leader-election-manager] got leader, starting workers...")
	for _, w := range m.workers {
		if err := w.Start(ctx); err != nil {
			slog.Error(fmt.Sprintf("[leader-election-manager] can not start worker: %s", err.Error()))
		}
	}
}

func (m *manager) onLeadingStop() {
	slog.Info("[leader-election-manager] lost leading, stopping workers...")
	for _, w := range m.workers {
		w.Stop()
	}
}

func (m *manager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}

	for _, w := range m.workers {
		w.Stop()
	}

	m.cancel()

LOOP:
	for {
		select {
		case <-m.done:
			slog.Info("leader election loop stopped")
			break LOOP
		case <-time.After(time.Second * 5):
			slog.Warn("leader election loop did not stop in 5s")
		}
	}
}
