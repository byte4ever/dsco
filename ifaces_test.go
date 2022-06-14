package dsco

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBindingAttempt_HasValue(t *testing.T) {
	t.Parallel()

	t.Run(
		"true case",
		func(t *testing.T) {
			t.Parallel()

			attempt := BindingAttempt{}

			require.False(t, attempt.HasValue())
		},
	)

	t.Run(
		"false case",
		func(t *testing.T) {
			t.Parallel()

			attempt := BindingAttempt{
				Value: reflect.ValueOf(10),
			}

			require.True(t, attempt.HasValue())
		},
	)
}
