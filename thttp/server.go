package thttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
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

	done chan struct{}
}

// NewGroup creates a new Group controlled by the given context
func NewGroup(ctx context.Context) *Group {
	g := new(Group)
	g.done = make(chan struct{})
	return g
}
func (g *Group) Spawn(task Task) {
	go g.runTask(g.ctx, 0, "", task)
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
	return task(ctx)
}

type contextKey int

const PanicCounterKey contextKey = iota

// Second parameter is the task ID. It is ignored because the only reason to
// pass it is to add it to the stack trace
func (g *Group) runTask(ctx context.Context, _ int64, name string, task Task) {
	err := RunTask(ctx, task)
	//	tlog.Get(ctx).Debug("Task finished", zap.Error(err))

	if err != nil {
		g.exit(err)
	}
}

func (g *Group) exit(err error) {
}

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
