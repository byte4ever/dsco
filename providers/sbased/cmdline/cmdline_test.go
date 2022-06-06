package cmdline

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/providers/sbased"
)

func TestProvide(t *testing.T) {
	type args struct {
		optionsLine []string
	}

	tests := []struct {
		name               string
		args               args
		want               *EntriesProvider
		wantErr            error
		invalidArgPosition int
	}{
		{
			name: "no options",
			args: args{
				optionsLine: []string{},
			},
			want: &EntriesProvider{
				values: nil,
			},
			wantErr: nil,
		},
		{
			name: "nil value",
			args: args{
				optionsLine: nil,
			},
			want: &EntriesProvider{
				values: nil,
			},
			wantErr: nil,
		},
		{
			name: "invalid format in first position",
			args: args{
				optionsLine: []string{"invalid_arg", "--arg1=value1", "--arg2=value2"},
			},
			want:               nil,
			wantErr:            ErrParamFormat,
			invalidArgPosition: 0,
		},
		{
			name: "invalid format in last position",
			args: args{
				optionsLine: []string{"--arg1=value1", "--arg2=value2", "invalid_arg"},
			},
			want:               nil,
			wantErr:            ErrParamFormat,
			invalidArgPosition: 2,
		},
		{
			name: "invalid format in middle position 1",
			args: args{
				optionsLine: []string{"--arg1=value1", "invalid_arg", "--arg2=value2"},
			},
			want:               nil,
			wantErr:            ErrParamFormat,
			invalidArgPosition: 1,
		},
		{
			name: "invalid format in middle position 2",
			args: args{
				optionsLine: []string{"--arg1=value1", "--asd-_asd=failure", "--arg2=value2"},
			},
			want:               nil,
			wantErr:            ErrParamFormat,
			invalidArgPosition: 1,
		},
		{
			name: "success single command line option",
			args: args{
				optionsLine: []string{"--arg1=value1"},
			},
			want: &EntriesProvider{
				values: sbased.Entries{
					"arg1": &sbased.Entry{
						ExternalKey: "--arg1",
						Value:       "value1",
					},
				},
			},
		},
		{
			name: "success multiple command line options",
			args: args{
				optionsLine: []string{
					"--arg1=value1",
					"--arg2=value2",
				},
			},
			want: &EntriesProvider{
				values: sbased.Entries{
					"arg1": &sbased.Entry{
						ExternalKey: "--arg1",
						Value:       "value1",
					},
					"arg2": &sbased.Entry{
						ExternalKey: "--arg2",
						Value:       "value2",
					},
				},
			},
		},
		{
			name: "duplicate params",
			args: args{
				optionsLine: []string{
					"--arg1=value1",
					"--arg1=value1x",
					"--arg2=value2",
				},
			},
			want:    nil,
			wantErr: ErrDuplicateParam,
		},
		{
			name: "with valid option",
			args: args{
				optionsLine: []string{
					"--arg1=value1",
					"--arg2=value2",
				},
			},
			want: &EntriesProvider{
				values: sbased.Entries{
					"arg1": &sbased.Entry{
						ExternalKey: "--arg1",
						Value:       "value1",
					},
					"arg2": &sbased.Entry{
						ExternalKey: "--arg2",
						Value:       "value2",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := NewEntriesProvider(tt.args.optionsLine)

				require.ErrorIsf(t, err, tt.wantErr, "NewEntriesProvider() error = %v, wantErr %v", err, tt.wantErr)
				require.Equalf(t, got, tt.want, "NewEntriesProvider() got = %v, want %v", got, tt.want)

				if err != nil {
					if errors.Is(err, ErrParamFormat) {
						require.ErrorContainsf(
							t,
							err,
							tt.args.optionsLine[tt.invalidArgPosition],
							"error message (%v) does not contains invalid arg content", err,
						)
						require.ErrorContainsf(
							t,
							err,
							fmt.Sprintf("#%d", tt.invalidArgPosition),
							"error message (%v) does not contains invalid arg position", err,
						)
						return
					}
				}
			},
		)
	}
}

func TestProvider_GetEntries(t *testing.T) {
	entries := sbased.Entries{
		"a1": &sbased.Entry{
			ExternalKey: "a",
			Value:       "b",
		},
		"b1": &sbased.Entry{
			ExternalKey: "a",
			Value:       "b",
		},
	}

	p := &EntriesProvider{
		values: entries,
	}

	require.Equal(t, entries, p.GetEntries())
}

func TestProvider_GetOrigin(t *testing.T) {
	require.Equal(t, dsco.Origin("cmdline"), (&EntriesProvider{}).GetOrigin())
}
