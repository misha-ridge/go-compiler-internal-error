package sorted

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapStringString(t *testing.T) {
	m := map[string]string{
		"q": "w",
		"e": "r",
		"t": "y",
		"u": "i",
	}

	assert.Equal(t, Keys(m), []string{"e", "q", "t", "u"})
}

func TestMapStringOther(t *testing.T) {
	m := map[string][]string{
		"q": {"w"},
		"e": {"r"},
		"t": {"y"},
		"u": {"i"},
	}

	assert.Equal(t, Keys(m), []string{"e", "q", "t", "u"})
}
