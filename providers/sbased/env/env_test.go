package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/providers/sbased"
)

func setEnv(t *testing.T, env map[string]string) {
	t.Helper()
	os.Clearenv()

	for k, v := range env {
		require.NoError(t, os.Setenv(k, v))
	}
}

func TestRegexpPrefix(t *testing.T) {
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

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				require.Equal(
					t,
					test.matches,
					getRePrefixed(test.prefix).FindStringSubmatch(test.str),
				)
			},
		)
	}
}

func TestRegexpKey(t *testing.T) {
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

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				require.Equal(t, test.matches, reSubKey.MatchString(test.str))
			},
		)
	}
}

func Test_RegexpPrefix(t *testing.T) {
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

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				require.Equal(t, test.match, rePrefix.MatchString(test.str))
			},
		)
	}
}

func TestProvide(t *testing.T) {
	type args struct {
		env    map[string]string
		prefix string
	}

	tests := []struct {
		name    string
		args    args
		want    *EntriesProvider
		wantErr []error
	}{
		{
			name: "",
			args: args{
				env:    nil,
				prefix: "PREFIX",
			},
			want: &EntriesProvider{
				prefix: "PREFIX",
			},
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
				prefix: "PREFIX",
				entries: sbased.Entries{
					"arg1": &sbased.Entry{
						ExternalKey: "PREFIX-ARG1",
						Value:       "val1",
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
				prefix: "PREFIX",
				entries: sbased.Entries{
					"arg1": &sbased.Entry{
						ExternalKey: "PREFIX-ARG1",
						Value:       "val1",
					},
					"arg1_v2": &sbased.Entry{
						ExternalKey: "PREFIX-ARG1_V2",
						Value:       "val2",
					},
				},
			},
		},
	}

	for _, tt := range tests {
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
	for _, prefix := range []string{
		"A A",
		"ad",
		"1A",
	} {
		t.Run(
			fmt.Sprintf("invalid prefix %s", prefix), func(t *testing.T) {
				p, errs := NewEntriesProvider(prefix)
				require.Len(t, errs, 1)
				require.ErrorIs(t, errs[0], ErrInvalidPrefix)
				require.ErrorContains(t, errs[0], prefix)
				require.Nil(t, p)
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
		entries: entries,
	}

	require.Equal(t, entries, p.GetEntries())
}

func TestProvider_GetOrigin(t *testing.T) {
	require.Equal(t, dsco.Origin("env"), (&EntriesProvider{}).GetOrigin())
}

func TestProvideKeySyntaxError(t *testing.T) {
	t.Run(
		"", func(t *testing.T) {
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

			p, errs := NewEntriesProvider("PREFIX")
			require.Nil(t, p)

			for i, err := range errs {
				require.ErrorIs(t, err, ErrInvalidKeyFormat)
				require.ErrorContains(t, err, keys[i])
			}
		},
	)
}
