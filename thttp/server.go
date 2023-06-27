package thttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
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
	done    chan struct{}
	closing bool
	err     error
}

// NewGroup creates a new Group controlled by the given context
func NewGroup(ctx context.Context) *Group {
	g := new(Group)
	g.done = make(chan struct{})
	return g
}
func (g *Group) Spawn(task Task) {
	g.mu.Lock()
	g.mu.Unlock()

	go g.runTask(g.ctx /*tlog.WithLogger(g.ctx, logger)*/, 0, "", 0, task)
}

// ErrPanic is the error type that occurs when a subtask panics
type ErrPanic struct {
	value any
	stack []byte
}

func (err ErrPanic) Error() string {
	return fmt.Sprintf("panic: %s", err.value)
}

// Unwrap returns the error passed to panic, or nil if panic was called with
// something other than an error
func (err ErrPanic) Unwrap() error {
	if e, ok := err.value.(error); ok {
		return e
	}
	return nil
}

// Value returns the value passed to panic
func (err ErrPanic) Value() any {
	return err.value
}

// Stack returns the panic stack trace
func (err ErrPanic) Stack() []byte {
	return err.stack
}

// RunTask executes the task in the current goroutine, recovering from panics.
// A panic is logged, reported to monitoring and returned as ErrPanic.
func RunTask(ctx context.Context, task Task) (err error) {
	defer func() {
		if p := recover(); p != nil {
			panicErr := ErrPanic{value: p, stack: debug.Stack()}
			err = panicErr
			//			tlog.Get(ctx).Error("Panic", zap.String("value", fmt.Sprint(p)), zap.ByteString("stack", panicErr.stack))
			if panicCounter := ctx.Value(PanicCounterKey); panicCounter != nil {
				panicCounter.(func())()
			}
		}
	}()
	return task(ctx)
}

type contextKey int

const PanicCounterKey contextKey = iota

// Second parameter is the task ID. It is ignored because the only reason to
// pass it is to add it to the stack trace
func (g *Group) runTask(ctx context.Context, _ int64, name string, onExit OnExit, task Task) {
	err := RunTask(ctx, task)
	//	tlog.Get(ctx).Debug("Task finished", zap.Error(err))

	g.mu.Lock()
	defer g.mu.Unlock()

	if err != nil {
		g.exit(err)
	} else if !g.closing {
		switch onExit {
		case Continue:
		case Exit:
			g.exit(nil)
		case Fail:
			g.exit(fmt.Errorf("task %q terminated unexpectedly", name))
		default:
			g.exit(fmt.Errorf("task %q: %v", name, onExit))
		}
	}
}

func (g *Group) exit(err error) {
}

// OnExit is an enumeration of exit handling modes. It specifies what should
// happen to the parent task if the subtask returns nil.
//
// Regardless of the chosen mode, if the subtask returns an error, it causes the
// parent task to shut down gracefully and return that error.
type OnExit int

const (
	// Continue means other subtasks of the parent task should continue to run.
	// Note that the parent task will return nil if its last remaining subtask
	// returns nil, even if Continue is specified.
	//
	// Use this mode for finite jobs that need to run once.
	Continue OnExit = iota

	// Exit means shut down the parent task gracefully.
	//
	// Use this mode for tasks that should be able to initiate graceful
	// shutdown, such as an HTTP server with a /quit endpoint that needs to
	// cause the process to exit.
	//
	// If any of other subtasks return an error, and it is not a (possibly
	// wrapped) context.Canceled, then the parent task will return the error.
	// Only first error from subtasks will be returned, the rest will be
	// discarded.
	//
	// If all other subtasks return nil or context.Canceled, the parent task
	// returns nil.
	Exit

	// Fail means shut down the parent task gracefully and return an error.
	//
	// Use this mode for subtasks that should never return unless their context
	// is closed.
	Fail
)

type SpawnFn func(task Task)
type Task func(ctx context.Context) error

// NewServer creates a Server
func NewServer(listener net.Listener, handler http.Handler) *Server {
	return &Server{}
}

func Run(ctx context.Context, start func(spawn SpawnFn) error) error {
	g := NewGroup(ctx)
	start(g.Spawn)
	return nil
}

// Run serves requests until the context is closed, then performs graceful
// shutdown for up to gracefulShutdownTimeout
func (s *Server) Run(ctx context.Context) error {
	return Run(ctx, func(spawn SpawnFn) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}
		return nil
	})
}

func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
