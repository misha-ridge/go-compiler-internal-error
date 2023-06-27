package time

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	now := Now()

	name, offset := now.Zone()
	assert.Equal(t, "UTC", name)
	assert.Equal(t, 0, offset)

	assert.Equal(t, "UTC", now.Location().String())
}
