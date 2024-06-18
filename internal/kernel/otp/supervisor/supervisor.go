package supervisor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/optclblast/drain-chain/internal/kernel/logger"
	"github.com/optclblast/drain-chain/internal/kernel/otp/application"
)

// A supervisor is a process which supervises other processes, which we refer to as child processes.
// The act of supervising a process includes three distinct responsibilities. The first one is to
// start child processes. Once a child process is running, the supervisor may restart a child
// process, either because it terminated abnormally or because a certain condition was reached.
// For example, a supervisor may restart all children if any child dies. Finally, a supervisor is
// also responsible for shutting down the child processes when the system is shutting down.
type Supervisor interface {
	// StartLink
	StartLink(args any) // todo replace with typed args

	// TODO: maybe replace application with ChildSpec type? But it seems irrelevant at the time
	Init(children []application.Application, strategy Strategy)
}

type Strategy interface {
	Handle(ctx context.Context, crashed application.Application) error
}

type defaultSupervisor struct {
	done     chan struct{}
	strategy Strategy
	log      *logger.Logger

	children []application.Application
}

func NewDefaultSupervisor(
	doneCh chan struct{},
	log *logger.Logger,
) Supervisor {
	return &defaultSupervisor{
		done: doneCh,
		log:  log,
	}
}

func (s *defaultSupervisor) StartLink(args any) {
	// TODO maybe here we need to register our supervisor in a some sort
	// of master-root-supervisor, as it made by calling Supervisor.start_link(__MODULE__, :ok, args)?
}

func (s *defaultSupervisor) Init(children []application.Application, strategy Strategy) {
	s.strategy = strategy
	s.children = children

	// start link with every child
	for _, child := range s.children {
		s.childLink(child)
	}
}

func (s *defaultSupervisor) childLink(child application.Application) {
	defer func() {
		if p := recover(); p != nil {
			s.log.Error(
				"supervisor.childLink: child paniced!",
				slog.Int("supervision_strategy", 0),
				slog.Any("panic", p),
				// TODO add child id name and some other info
			)

			if err := s.strategy.Handle(context.Background(), child); err != nil {
				s.log.Error(
					"supervisor.childLink: error handle one of child's crash with strategy",
					logger.Err(err),
				)
			}
		}
	}()

	child.Start()
}

type OneForAllStrategy struct {
	Children []application.Application

	m          sync.RWMutex
	restarting bool
}

func (s *OneForAllStrategy) Handle(ctx context.Context, crashed application.Application) error {
	s.m.RLock()
	restarting := s.restarting
	s.m.RUnlock()

	if restarting {
		return nil
	}

	s.m.Lock()
	s.restarting = true
	s.m.Unlock()

	done := make(chan struct{}, 1)

	go func() {
		defer close(done)

		for _, child := range s.Children {
			select {
			case <-ctx.Done():
				break
			default:
				child.Reload(ctx)
			}
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
