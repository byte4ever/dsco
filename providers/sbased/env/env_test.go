package env

import (
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

func Test_Regexp(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		matches []string
	}{
		{
			name: "",
			str:  "really wrong",
		},
		{
			name: "",
			str:  "A-B-C",
		},
		{
			name: "",
			str:  "--A-B-C",
		},
		{
			name: "",
			str:  "PRE_FIX-B-C_D=abcdef",
		},
		{
			name: "",
			str:  "1PREFIX-B-C_D=abcdef",
		},
		{
			name: "",
			str:  "PREFIX.B-C_D=abcdef",
		},
		{
			name: "",
			str:  "A1A-1B-C_D=abcdef",
		},
		{
			name: "",
			str:  "A1A-B-1C_D=abcdef",
		},
		{
			name:    "",
			str:     "A-B-C=abcdef",
			matches: []string{"A-B-C=abcdef", "A", "B-C", "abcdef"},
		},
		{
			name:    "",
			str:     "A-B-C_D=abcdef",
			matches: []string{"A-B-C_D=abcdef", "A", "B-C_D", "abcdef"},
		},
		{
			name:    "",
			str:     "A1-B-C_D=abcdef",
			matches: []string{"A1-B-C_D=abcdef", "A1", "B-C_D", "abcdef"},
		},
		{
			name:    "",
			str:     "A1A-B-C_D=abcdef",
			matches: []string{"A1A-B-C_D=abcdef", "A1A", "B-C_D", "abcdef"},
		},
		{
			name:    "",
			str:     "PREFIX-B-C_D=abcdef",
			matches: []string{"PREFIX-B-C_D=abcdef", "PREFIX", "B-C_D", "abcdef"},
		},
		{
			name:    "",
			str:     "A1A-B-C_D=abc=def",
			matches: []string{"A1A-B-C_D=abc=def", "A1A", "B-C_D", "abc=def"},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				require.Equal(t, test.matches, re.FindStringSubmatch(test.str))
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
		want    *Provider
		wantErr error
	}{
		{
			name: "",
			args: args{
				env:    nil,
				prefix: "A A",
			},
			want:    nil,
			wantErr: ErrInvalidPrefix,
		},
		{
			name: "",
			args: args{
				env:    nil,
				prefix: "PREFIX",
			},
			want: &Provider{
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
			want: &Provider{
				prefix: "PREFIX",
				entries: sbased.StrEntries{
					"arg1": &sbased.StrEntry{
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
			want: &Provider{
				prefix: "PREFIX",
				entries: sbased.StrEntries{
					"arg1": &sbased.StrEntry{
						ExternalKey: "PREFIX-ARG1",
						Value:       "val1",
					},
					"arg1_v2": &sbased.StrEntry{
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
				got, err := Provide(tt.args.prefix)

				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tt.want, got)
			},
		)
	}
}

func TestProvider_GetEntries(t *testing.T) {
	entries := sbased.StrEntries{
		"a1": &sbased.StrEntry{
			ExternalKey: "a",
			Value:       "b",
		},
		"b1": &sbased.StrEntry{
			ExternalKey: "a",
			Value:       "b",
		},
	}

	p := &Provider{
		entries: entries,
	}

	require.Equal(t, entries, p.GetEntries())
}

func TestProvider_GetOrigin(t *testing.T) {
	require.Equal(t, dsco.Origin("env"), (&Provider{}).GetOrigin())
}
