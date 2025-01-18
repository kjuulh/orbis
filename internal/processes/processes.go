package processes

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type Process interface {
	Start(ctx context.Context) error
}

type SetupProcesser interface {
	Setup(ctx context.Context) error
}

type CloseProcesser interface {
	Close(ctx context.Context) error
}

type App struct {
	logger    *slog.Logger
	processes []Process

	externalDone chan bool
}

func NewApp(logger *slog.Logger) *App {
	return &App{
		logger:       logger,
		processes:    make([]Process, 0),
		externalDone: make(chan bool),
	}
}

func (a *App) Add(p Process) *App {
	a.processes = append(a.processes, p)

	return a
}

func (a *App) WithSignal(stop func()) *App {
	go func() {
		stop()

		a.externalDone <- true
	}()

	return a
}

func (a *App) WithCtrlC() *App {
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		<-stop

		a.externalDone <- true
	}()

	return a
}

func (a *App) Execute(ctx context.Context) error {
	a.logger.InfoContext(ctx, "starting processor")
	if err := a.setupProcesses(ctx); err != nil {
		return err
	}

	processes, err := a.startProcesses(ctx)
	if err != nil {
		return nil
	}

	processErr := processes.wait(ctx, a.externalDone)

	if err := a.closeProcesses(ctx, processes); err != nil {
		if processErr != nil {
			return processErr
		}

		return err
	}

	if processErr != nil {
		return processErr
	}

	return nil
}

func (a *App) closeProcesses(ctx context.Context, processes *processStatus) error {
	waitClose, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	closeErrs := make(chan error)

	go func() {
		errgrp, ctx := errgroup.WithContext(waitClose)
		for _, closeProcessor := range a.processes {
			if close, ok := closeProcessor.(CloseProcesser); ok {
				errgrp.Go(func() error {
					a.logger.InfoContext(ctx, "closing processor")
					return close.Close(ctx)
				})

			}
		}

		closeErrs <- errgrp.Wait()
	}()

	select {
	case <-waitClose.Done():
		return nil
	case <-closeErrs:
		return nil
	case _, closed := <-processes.errs:
		if closed {
			return nil
		}
	}

	for _, closeHandle := range processes.processHandles {
		closeHandle()
	}

	return nil
}

type processStatus struct {
	errs           chan error
	processHandles []context.CancelFunc
}

func (p *processStatus) wait(_ context.Context, externalDone chan bool) error {
	select {
	case err := <-p.errs:
		return err
	case <-externalDone:
		return nil
	}
}

func (a *App) startProcesses(ctx context.Context) (*processStatus, any) {
	status := &processStatus{
		errs:           make(chan error, len(a.processes)),
		processHandles: make([]context.CancelFunc, 0),
	}

	for _, process := range a.processes {
		processCtx, cancelFunc := context.WithCancel(ctx)

		status.processHandles = append(status.processHandles, cancelFunc)

		go func(ctx context.Context, process Process) {
			a.logger.DebugContext(ctx, "starting process")

			err := process.Start(ctx)

			if err != nil {
				a.logger.WarnContext(ctx, "process finished with error", "error", err)

			} else {
				a.logger.DebugContext(ctx, "process finished gracefully")

			}

			status.errs <- err
		}(processCtx, process)
	}

	return status, nil
}

func (a *App) setupProcesses(ctx context.Context) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	errgrp, ctx := errgroup.WithContext(ctxWithDeadline)
	for _, setupProcessor := range a.processes {
		if setup, ok := setupProcessor.(SetupProcesser); ok {
			errgrp.Go(func() error {
				a.logger.InfoContext(ctx, "setting up processor")
				return setup.Setup(ctx)
			})
		}
	}

	return errgrp.Wait()
}
