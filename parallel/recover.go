package parallel

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/misha-ridge/x/tlog"
	"go.uber.org/zap"
)

type contextKey int

const (
	// PanicCounterKey is a context key for panic counter metric registration.
	//
	// Register a func() under this key to have it called on every panic captured by parallel.
	PanicCounterKey contextKey = iota
)

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
			tlog.Get(ctx).Error("Panic", zap.String("value", fmt.Sprint(p)), zap.ByteString("stack", panicErr.stack))
			if panicCounter := ctx.Value(PanicCounterKey); panicCounter != nil {
				panicCounter.(func())()
			}
		}
	}()
	return task(ctx)
}
