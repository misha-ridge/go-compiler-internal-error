package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFakeClock(t *testing.T) {
	assert.True(t, IsRealClock())

	now := Now()

	fakeClock := SetFakeClockForTesting(t)
	assert.False(t, IsRealClock())

	fakeNow := Now()

	fakeClock.Advance(time.Minute)
	assert.Equal(t, time.Minute, Now().Sub(fakeNow))

	SetRealClock()
	time.Sleep(10 * time.Millisecond)
	now2 := Now()
	assert.True(t, now2.After(now))
	assert.True(t, IsRealClock())
}
