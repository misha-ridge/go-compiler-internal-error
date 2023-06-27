package time

import (
	"testing"

	"github.com/DPJacques/clockwork"
)

// Clock provides an interface that packages can use instead of directly
// using the time module, so that chronology-related behavior can be tested
type Clock = clockwork.Clock

// Timer reexports clockwork Timer type
type Timer = clockwork.Timer

// Ticker reexports clockwork Ticker type
type Ticker = clockwork.Ticker

// You have arrived here in a search of a source of data race found by
// go test -race. You see a variable `clock` accessed from several goroutines
// without a mutex. You say "a-ha!" and prepare a PR to add a lock around it.
//
// Don't.
//
// The bug is in the test, not in the library.
//
// Tests ought to change the clock to the fake one before the test starts, then
// execute the test itself, then reap any stray goroutines, and then revert the
// clock to the real one.
//
// If you see a data race between accessing a clock in goroutine run by test A
// and SetFakeClockForTesting[At] run by the prologue of test B, then it means A
// started the goroutine, but hasn't terminated it when it finished. Test B was
// started, it changed the clock to fake one, and this race was caught by the
// data race checker.
//
// Stray goroutines are a problem: if they are not strictly read-only (and why
// would they?), they affect tests that run after them. Find them, and tidy them
// up before finishing the test.
var clock = realClock

var realClock = clockwork.NewRealClock()

// SetFakeClockForTestingAt replaces the system clock a with a fake clock at the
// given time for testing. Real clock is reinstalled at the end of the test.
func SetFakeClockForTestingAt(t *testing.T, time Time) clockwork.FakeClock {
	c := clockwork.NewFakeClockAt(time)
	clock = c
	t.Cleanup(SetRealClock)
	return c
}

// SetFakeClockForTesting replaces the system clock with a fake clock for
// testing. Real clock is reinstalled at the end of the test.
func SetFakeClockForTesting(t *testing.T) clockwork.FakeClock {
	c := clockwork.NewFakeClock()
	clock = c
	t.Cleanup(SetRealClock)
	return c
}

// SetRealClock returns the system clock to the real clock
func SetRealClock() {
	clock = realClock
}

// IsRealClock returns true if the single instance of the clock represent real time
func IsRealClock() bool {
	return clock == realClock
}

// RealClock returns the real clock
func RealClock() Clock {
	return realClock
}
