package dsco

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoundingAttempt_HasValue(t *testing.T) {
	t.Parallel()

	t.Run(
		"true case",
		func(t *testing.T) {
			t.Parallel()

			attempt := BoundingAttempt{}

			require.False(t, attempt.HasValue())
		},
	)

	t.Run(
		"false case",
		func(t *testing.T) {
			t.Parallel()

			attempt := BoundingAttempt{
				Value: reflect.ValueOf(10),
			}

			require.True(t, attempt.HasValue())
		},
	)
}
