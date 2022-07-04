package dsco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_convert(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single word",
			args: args{
				s: "Hello",
			},
			want: "hello",
		},
		{
			name: "two words",
			args: args{
				s: "HelloWorld",
			},
			want: "hello_world",
		},
		{
			name: "two words",
			args: args{
				s: "Hello.World",
			},
			want: "hello-world",
		},
		{
			name: "mixing",
			args: args{
				s: "HelloWorld1.HelloWorld2",
			},
			want: "hello_world1-hello_world2",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				require.Equal(t, tt.want, convert(tt.args.s))
			},
		)
	}
}
