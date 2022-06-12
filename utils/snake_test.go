package utils

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	t.Parallel()

	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "T1",
			args: args{
				str: "AbaloneToto",
			},
			want: "abalone_toto",
		},
		{
			name: "T2",
			args: args{
				str: "AbaloneTotoPolo",
			},
			want: "abalone_toto_polo",
		},
		{
			name: "T3",
			args: args{
				str: "Abalone",
			},
			want: "abalone",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				if got := toSnakeCase(tt.args.str); got != tt.want {
					t.Errorf("toSnakeCase() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
