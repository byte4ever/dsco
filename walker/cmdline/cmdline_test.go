package cmdline

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/walker/svalues"
)

func TestProvide(t *testing.T) {
	t.Parallel()

	type args struct {
		optionsLine []string
	}

	const (
		arg1        = "--arg1=value1"
		arg2        = "--arg2=value2"
		arg3        = "--arg3=value3"
		invalidArg1 = "invalid_arg1"
		invalidArg2 = "--asd-_asd=failure"
	)

	tests := []struct {
		wantErrorIs        []error
		wantErrorPositions []int
		wantErrorContains  [][]string
		want               *EntriesProvider
		name               string
		args               args
	}{
		{
			name: "no options",
			args: args{
				optionsLine: []string{},
			},
			want:        &EntriesProvider{},
			wantErrorIs: nil,
		},
		{
			name: "nil value",
			args: args{
				optionsLine: nil,
			},
			want: &EntriesProvider{},
		},
		{
			name: "invalid format in first position",
			args: args{
				optionsLine: []string{
					invalidArg1,
					arg1,
					arg2,
				},
			},
			want:               nil,
			wantErrorIs:        []error{ErrInvalidFormat},
			wantErrorPositions: []int{1},
			wantErrorContains:  [][]string{{invalidArg1}},
		},
		{
			name: "invalid format in last position",
			args: args{
				optionsLine: []string{
					arg1,
					arg2,
					invalidArg1,
				},
			},
			want:               nil,
			wantErrorIs:        []error{ErrInvalidFormat},
			wantErrorPositions: []int{3},
			wantErrorContains:  [][]string{{invalidArg1}},
		},
		{
			name: "invalid format in middle position 1",
			args: args{
				optionsLine: []string{
					arg1,
					arg2,
					invalidArg1,
					arg3,
				},
			},
			want:               nil,
			wantErrorIs:        []error{ErrInvalidFormat},
			wantErrorPositions: []int{3},
			wantErrorContains:  [][]string{{invalidArg1}},
		},
		{
			name: "invalid format in middle position 2",
			args: args{
				optionsLine: []string{
					arg1,
					invalidArg2,
					arg2,
				},
			},
			want:               nil,
			wantErrorIs:        []error{ErrInvalidFormat},
			wantErrorPositions: []int{2},
			wantErrorContains:  [][]string{{invalidArg2}},
		},
		{
			name: "success single command line option",
			args: args{
				optionsLine: []string{arg1},
			},
			want: &EntriesProvider{
				stringValues: svalues.StringValues{
					"arg1": {
						Location: "cmdline[--arg1]",
						Value:    "value1",
					},
				},
			},
		},
		{
			name: "success multiple command line options",
			args: args{
				optionsLine: []string{
					arg1,
					arg2,
				},
			},
			want: &EntriesProvider{
				stringValues: svalues.StringValues{
					"arg1": {
						Location: "cmdline[--arg1]",
						Value:    "value1",
					},
					"arg2": {
						Location: "cmdline[--arg2]",
						Value:    "value2",
					},
				},
			},
		},
		{
			name: "duplicate params",
			args: args{
				optionsLine: []string{
					arg1,
					"--arg1=value1x",
					arg2,
				},
			},
			want:               nil,
			wantErrorIs:        []error{ErrDuplicateParam},
			wantErrorPositions: []int{2},
			wantErrorContains:  [][]string{{"arg1", "#1", "previous"}},
		},
		{
			name: "mixed errors",
			args: args{
				optionsLine: []string{
					invalidArg1,
					arg1,
					"--arg1=value1x",
					arg2,
					arg3,
					invalidArg2,
				},
			},
			wantErrorIs: []error{
				ErrInvalidFormat,
				ErrDuplicateParam,
				ErrInvalidFormat,
			},
			wantErrorPositions: []int{1, 3, 6},
			wantErrorContains: [][]string{
				{invalidArg1},
				{"arg1", "#2", "previous"},
				{invalidArg2},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				// test precondition
				require.Equal(
					t,
					len(tt.wantErrorPositions),
					len(tt.wantErrorContains),
				)
				require.Equal(
					t,
					len(tt.wantErrorContains),
					len(tt.wantErrorIs),
				)

				provider, err := NewEntriesProvider(tt.args.optionsLine)
				if len(tt.wantErrorPositions) > 0 {
					require.Nil(t, provider)
					// check for error
					var pe *ParamError
					require.ErrorAs(t, err, &pe)

					require.Equal(t, tt.wantErrorPositions, pe.Positions)

					for i, subError := range tt.wantErrorIs {
						require.ErrorIsf(
							t,
							pe.Errs[i],
							subError,
							"error %d",
							i,
						)
					}

					for i, subErrorContains := range tt.wantErrorContains {
						for _, contain := range subErrorContains {
							require.ErrorContainsf(
								t,
								pe.Errs[i],
								contain,
								"error %d",
								i,
							)
						}
					}

					return
				}

				require.NoError(t, err)
				require.NotNil(t, provider)
				require.EqualValues(t, tt.want, provider)

			},
		)
	}
}

func TestProvider_GetStringValues(t *testing.T) {
	t.Parallel()

	entries := svalues.StringValues{
		"k1": {
			Location: "l1",
			Value:    "v1",
		},
		"k2": {
			Location: "l2",
			Value:    "v2",
		},
	}

	p := &EntriesProvider{
		stringValues: entries,
	}

	require.Equal(t, entries, p.GetStringValues())
}
