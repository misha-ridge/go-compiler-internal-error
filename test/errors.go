package test

import (
	"errors"
	"testing"

	"github.com/misha-ridge/x/apierror"
	"github.com/stretchr/testify/require"
)

// RequireStatus requires an error to be of the type apierror and have the expected status code
func RequireStatus(t *testing.T, err error, expected int) {
	var e apierror.Error
	require.Error(t, err)
	require.True(t, errors.As(err, &e), "error is not an api error")
	require.Equal(t, expected, e.HTTPStatus())
}
