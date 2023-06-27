package thttp

import (
	"context"
	"errors"
	"fmt"
	"github.com/misha-ridge/x/tlog"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/misha-ridge/x/parallel"
)

// Server wraps an HTTP server
type Server struct {
}

type Group struct {
	// group controls lifetimes of its members via this context.
	// Group is decoupled from the call stack, so it has to carry
	// context in the struct, not in a parameter
	ctx    context.Context //nolint:containedctx
	cancel context.CancelFunc

	mu      sync.Mutex
	running int
	done    chan struct{}
	closing bool
	err     error
}

// NewGroup creates a new Group controlled by the given context
func NewGroup(ctx context.Context) *Group {
	g := new(Group)
	g.ctx, g.cancel = context.WithCancel(ctx)
	g.done = make(chan struct{})
	close(g.done)
	return g
}

var nextTaskID int64 = 0x0bace1d000000000

func (g *Group) Spawn(name string, onExit parallel.OnExit, task parallel.Task) {
	id := atomic.AddInt64(&nextTaskID, 1)

	g.mu.Lock()
	if g.running == 0 {
		g.done = make(chan struct{})
	}
	g.running++
	g.mu.Unlock()

	logger := tlog.Get(g.ctx).Named(name)
	logger.Debug("Task spawned", zap.String("id", fmt.Sprintf("%x", id)), zap.Stringer("onExit", onExit))

	go g.runTask(tlog.WithLogger(g.ctx, logger), id, name, onExit, task)
}

// Second parameter is the task ID. It is ignored because the only reason to
// pass it is to add it to the stack trace
func (g *Group) runTask(ctx context.Context, _ int64, name string, onExit parallel.OnExit, task parallel.Task) {
	err := parallel.RunTask(ctx, task)
	tlog.Get(ctx).Debug("Task finished", zap.Error(err))

	g.mu.Lock()
	defer g.mu.Unlock()

	if err != nil {
		g.exit(err)
	} else if !g.closing {
		switch onExit {
		case parallel.Continue:
		case parallel.Exit:
			g.exit(nil)
		case parallel.Fail:
			g.exit(fmt.Errorf("task %q terminated unexpectedly", name))
		default:
			g.exit(fmt.Errorf("task %q: %v", name, onExit))
		}
	}

	g.running--
	if g.running == 0 {
		close(g.done)
	}
}

func (g *Group) exit(err error) {
	// Cancellations during shutdown are fine
	if g.closing && errors.Is(err, context.Canceled) {
		return
	}
	if g.err == nil {
		g.err = err
	}
	if !g.closing {
		g.closing = true
		g.cancel()
	}
}

type SpawnFn func(name string, onExit parallel.OnExit, task parallel.Task)

// NewServer creates a Server
func NewServer(listener net.Listener, handler http.Handler) *Server {
	return &Server{}
}

func Run(ctx context.Context, start func(ctx context.Context, spawn SpawnFn) error) error {
	g := NewGroup(ctx)
	start(nil, g.Spawn)
	return nil
}

// Run serves requests until the context is closed, then performs graceful
// shutdown for up to gracefulShutdownTimeout
func (s *Server) Run(ctx context.Context) error {
	return Run(ctx, func(ctx context.Context, spawn SpawnFn) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}
		return nil
	})
}

func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
