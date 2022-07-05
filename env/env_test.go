package env

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/svalue"
)

func setEnv(t *testing.T, env map[string]string) {
	t.Helper()

	for k, v := range env {
		t.Setenv(k, v)
	}
}

func TestRegexpPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		str     string
		prefix  string
		matches []string
	}{
		{
			name:   "",
			str:    "really wrong",
			prefix: "PREFIX",
		},
		{
			name:   "",
			str:    "A-B-C",
			prefix: "A",
		},
		{
			name:   "",
			str:    "--A-B-C",
			prefix: "PREFIX",
		},
		{
			name:   "",
			str:    "PRE_FIX-B-C_D=abcdef",
			prefix: "PREFIX",
		},
		{
			name:    "",
			str:     "PREFIX-1B-C_D=abcdef",
			matches: []string{"PREFIX-1B-C_D=abcdef", "-1B-C_D", "abcdef"},
			prefix:  "PREFIX",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				require.Equal(
					t,
					tt.matches,
					getRePrefixed(tt.prefix).FindStringSubmatch(tt.str),
				)
			},
		)
	}
}

func TestRegexpKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		str     string
		matches bool
	}{
		{
			name: "",
			str:  "really wrong",
		},
		{
			name:    "",
			str:     "-A-B-C",
			matches: true,
		},
		{
			name: "",
			str:  "--A-B-C",
		},
		{
			name: "",
			str:  "-1B-C_D=abcdef",
		},
		{
			name: "",
			str:  "-B-1C_D=abcdef",
		},
		{
			name: "",
			str:  "-.B-C_D=abcdef",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				require.Equal(t, tt.matches, reSubKey.MatchString(tt.str))
			},
		)
	}
}

func Test_RegexpPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		str   string
		match bool
	}{
		{
			name: "",
			str:  "really wrong",
		},
		{
			name: "",
			str:  "as",
		},
		{
			name: "",
			str:  "A B",
		},
		{
			name: "",
			str:  "1AC",
		},
		{
			name:  "",
			str:   "AC",
			match: true,
		},
		{
			name:  "",
			str:   "A1C",
			match: true,
		},
		{
			name:  "",
			str:   "AC1",
			match: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				require.Equal(t, tt.match, rePrefix.MatchString(tt.str))
			},
		)
	}
}

func TestProvide(t *testing.T) { //nolint:paralleltest // using env variable
	type args struct {
		env    map[string]string
		prefix string
	}

	tests := []struct {
		wantErr error
		want    *EntriesProvider
		args    args
		name    string
	}{
		{
			name: "",
			args: args{
				env:    nil,
				prefix: "PREFIX",
			},
			want: &EntriesProvider{},
		},
		{
			name: "",
			args: args{
				env: map[string]string{
					"ASAS":        "234",
					"PREFIX-ARG1": "val1",
				},
				prefix: "PREFIX",
			},
			want: &EntriesProvider{
				stringValues: svalue.Values{
					"arg1": {
						Location: "env[PREFIX-ARG1]",
						Value:    "val1",
					},
				},
			},
		},
		{
			name: "",
			args: args{
				env: map[string]string{
					"ASAS":           "234",
					"PREFIX-ARG1":    "val1",
					"PREFIX-ARG1_V2": "val2",
				},
				prefix: "PREFIX",
			},
			want: &EntriesProvider{
				stringValues: svalue.Values{
					"arg1": {
						Location: "env[PREFIX-ARG1]",
						Value:    "val1",
					},
					"arg1_v2": {
						Location: "env[PREFIX-ARG1_V2]",
						Value:    "val2",
					},
				},
			},
		},
	}

	for _, tt := range tests { //nolint:paralleltest // using setenv
		t.Run(
			tt.name, func(t *testing.T) {
				setEnv(t, tt.args.env)
				got, errs := NewEntriesProvider(tt.args.prefix)

				require.Equal(t, tt.wantErr, errs)
				require.Equal(t, tt.want, got)
			},
		)
	}
}

func TestProvide_InvalidPrefix(t *testing.T) {
	t.Parallel()

	//nolint:paralleltest // buggy linter
	for _, prefix := range []string{
		"A A",
		"ad",
		"1A",
	} {
		prefix := prefix
		t.Run(
			fmt.Sprintf("invalid prefix %s", prefix),
			func(t *testing.T) {
				t.Parallel()

				p, err := NewEntriesProvider(prefix)
				require.ErrorIs(t, err, ErrInvalidPrefix)
				require.ErrorContains(t, err, prefix)
				require.Nil(t, p)
			},
		)
	}
}

func TestProvider_GetStringValues(t *testing.T) {
	t.Parallel()

	stringValues := svalue.Values{
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
		stringValues: stringValues,
	}

	require.Equal(t, stringValues, p.GetStringValues())
}

func TestProvideKeySyntaxError(t *testing.T) {
	t.Parallel()

	t.Run(
		"",
		func(t *testing.T) {
			keys := []string{
				"PREFIX-1A-B",
				"PREFIX-A-1B",
				"PREFIX-A-_B",
				"PREFIX-AS,DD",
				"PREFIX-a-b",
				"PREFIXASDD",
			}
			envMap := make(map[string]string, len(keys))

			for i, key := range keys {
				envMap[key] = fmt.Sprintf("Value%d", i)
			}

			setEnv(
				t, envMap,
			)

			p, err := NewEntriesProvider("PREFIX")
			require.Nil(t, p)

			require.ErrorIs(t, err, ErrAmbiguousKeys)
			for _, k := range keys {
				require.ErrorContains(t, err, k)
			}
		},
	)
}
