package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FormatIndexSequence(t *testing.T) {
	t.Parallel()

	type args struct {
		indexes []int
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "single idx",
			args: args{
				indexes: []int{123},
			},
			want: "#123",
		},
		{
			name: "2 values",
			args: args{
				indexes: []int{123, 4000},
			},
			want: "#123 and #4000",
		},
		{
			name: "3 values",
			args: args{
				indexes: []int{123, 4000, 233},
			},
			want: "#123, #4000 and #233",
		},
		{
			name: "many values",
			args: args{
				indexes: []int{1, 2, 3, 4},
			},
			want: "#1, #2, #3 and #4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				got := FormatIndexSequence(tt.args.indexes)
				require.Equal(
					t,
					tt.want,
					got,
					"FormatIndexSequence() = %v, want %v",
					got,
					tt.want,
				)
			},
		)
	}
}

func Test_formatIndexSequence_panics(t *testing.T) {
	t.Parallel()
	require.Panics(
		t, func() {
			FormatIndexSequence(nil)
		},
	)
	require.Panics(
		t, func() {
			FormatIndexSequence([]int{})
		},
	)
}

func Test_FormatStringSequence(t *testing.T) {
	t.Parallel()

	type args struct {
		values []string
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "single value",
			args: args{
				values: []string{"hello"},
			},
			want: `"hello"`,
		},
		{
			name: "2 values",
			args: args{
				values: []string{"hello", "world"},
			},
			want: `"hello" and "world"`,
		},
		{
			name: "3 values",
			args: args{
				values: []string{"hello", "my", "world"},
			},
			want: `"hello", "my" and "world"`,
		},
		{
			name: "many values",
			args: args{
				values: []string{"hello", "my", "nice", "world"},
			},
			want: `"hello", "my", "nice" and "world"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				got := FormatStringSequence(tt.args.values)
				require.Equal(
					t,
					tt.want,
					got,
					"FormatIndexSequence() = %v, want %v",
					got,
					tt.want,
				)
			},
		)
	}
}

func Test_FormatStringSequence_panics(t *testing.T) {
	t.Parallel()

	require.Panics(
		t, func() {
			FormatStringSequence(nil)
		},
	)
	require.Panics(
		t, func() {
			FormatStringSequence([]string{})
		},
	)
}
