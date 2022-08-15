package managers

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// External for testing
var stopChan chan os.Signal

type ShutdownManager interface {
	Start(startFn func() error, endFn func() error)
	Wait() error
	Cancel()
}

type shutdownManager struct {
	g      *errgroup.Group
	ctx    context.Context
	cancel func()
}

func NewShutdownManager() ShutdownManager {
	ctx, cancel := context.WithCancel(context.Background())
	g, gCtx := errgroup.WithContext(ctx)

	stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		cancel()
	}()

	return &shutdownManager{
		g:      g,
		ctx:    gCtx,
		cancel: cancel,
	}
}

func (m *shutdownManager) Start(startFn func() error, endFn func() error) {
	m.g.Go(func() error {
		return startFn()
	})
	m.g.Go(func() error {
		<-m.ctx.Done()
		return endFn()
	})
}

func (m *shutdownManager) Wait() error {
	return m.g.Wait()
}

func (m *shutdownManager) Cancel() {
	m.cancel()
}
