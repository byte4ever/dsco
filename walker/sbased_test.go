package walker

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/svalues"
)

type Sub struct {
	FirstName    *string
	LastName     *string
	TrainingTime *time.Duration
	T            *time.Time
	B            *bool
}

type Root struct {
	A    *int
	B    *float64
	T    *time.Time
	Z    *Sub
	NaNa *int
	L    []string
}

func TestNewStringBased(t *testing.T) {
	t.Parallel()

	pp := &Root{}

	provider := NewMockStringValuesProvider(t)
	provider.
		On("GetStringValues").
		Return(
			svalues.StringValues{
				"a": {
					Location: "loc1",
					Value:    "1234",
				},
				"z-last_name": {
					Location: "loc2",
					Value:    "MARTIN",
				},
				"z-first_name": {
					Location: "loc3",
					Value:    "Laurent",
				},
				"z-training_time": {
					Location: "loc4",
					Value:    "123s",
				},
				"GrosBolos": {
					Location: "loc[gros bolos]",
					Value:    "123s",
				},
			},
		).Once()

	builder, err := NewStringBasedBuilder(
		provider,
	)

	require.NoError(t, err)
	require.NotNil(t, builder)

	base, errs := builder.GetBaseFor(pp)

	for i, err := range errs {
		fmt.Println("---", i, err)
	}

	for i, e := range base {
		fmt.Println(
			"====", i, e.path, e.location,
			(*e.value).Elem().Interface(),
		)
	}

}
