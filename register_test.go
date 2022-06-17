package dsco

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	id1 "github.com/byte4ever/dsco/dummy/d1"
)

func Test_longTypeName(t *testing.T) {
	t.Parallel()

	type args struct {
		t reflect.Type
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "built in",
			args: args{
				t: reflect.TypeOf(0),
			},
			want: "int",
		},
		{
			name: "built in pointer",
			args: args{
				t: reflect.TypeOf(R(0)),
			},
			want: "*int",
		},
		{
			name: "go pkg",
			args: args{
				t: reflect.TypeOf(time.Time{}),
			},
			want: "time/time.Time",
		},
		{
			name: "go pkg pointer",
			args: args{
				t: reflect.TypeOf(&time.Time{}),
			},
			want: "*time/time.Time",
		},
		{
			name: "internal pkg",
			args: args{
				t: reflect.TypeOf(id1.T1{}),
			},
			want: "github.com/byte4ever/dsco/dummy/d1/id1.T1",
		},
		{
			name: "internal pkg pointer",
			args: args{
				t: reflect.TypeOf(&id1.T1{}),
			},
			want: "*github.com/byte4ever/dsco/dummy/d1/id1.T1",
		},
		{
			name: "external pkg",
			args: args{
				t: reflect.TypeOf(require.Assertions{}),
			},
			want: "github.com/stretchr/testify/require/require.Assertions",
		},
		{
			name: "external pkg pointer",
			args: args{
				t: reflect.TypeOf(&require.Assertions{}),
			},
			want: "*github.com/stretchr/testify/require/require.Assertions",
		},
		{
			name: "dummy",
			args: args{
				t: reflect.TypeOf(
					struct {
						X int
						Y float64
					}{},
				),
			},
			want: "struct { X int; Y float64 }",
		},
		{
			name: "dummy pointer",
			args: args{
				t: reflect.TypeOf(
					&struct {
						X int
						Y float64
					}{},
				),
			},
			want: "*struct { X int; Y float64 }",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				require.Equal(t, tt.want, longTypeName(tt.args.t))
			},
		)
	}
}

func Test_typeIsRegistered(t *testing.T) {
	t.Parallel()

	type args struct {
		v any
	}

	tests := []struct {
		args args
		name string
		want bool
	}{
		{
			name: "int",
			args: args{
				v: R(0),
			},
			want: true,
		},
		{
			name: "int8",
			args: args{
				v: R(int8(0)),
			},
			want: true,
		},
		{
			name: "int16",
			args: args{
				v: R(int16(0)),
			},
			want: true,
		},
		{
			name: "int32",
			args: args{
				v: R(int32(0)),
			},
			want: true,
		},
		{
			name: "int64",
			args: args{
				v: R(int64(0)),
			},
			want: true,
		},
		{
			name: "uint",
			args: args{
				v: R(uint(0)),
			},
			want: true,
		},
		{
			name: "uint8",
			args: args{
				v: R(uint8(0)),
			},
			want: true,
		},
		{
			name: "uint16",
			args: args{
				v: R(uint16(0)),
			},
			want: true,
		},
		{
			name: "uint32",
			args: args{
				v: R(uint32(0)),
			},
			want: true,
		},
		{
			name: "uint64",
			args: args{
				v: R(uint64(0)),
			},
			want: true,
		},
		{
			name: "float32",
			args: args{
				v: R(float32(0)),
			},
			want: true,
		},
		{
			name: "float64",
			args: args{
				v: R(float64(0)),
			},
			want: true,
		},
		{
			name: "string",
			args: args{
				v: R(""),
			},
			want: true,
		},
		{
			name: "time.Time",
			args: args{
				v: R(&time.Time{}),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				require.Equal(
					t,
					tt.want,
					TypeIsRegistered(reflect.TypeOf(tt.args.v)),
				)
			},
		)
	}
}

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run(
		"",
		func(t *testing.T) {
			t.Parallel()

			typeToRegister := &id1.T1{}

			// should not be registered at startup
			require.False(
				t,
				TypeIsRegistered(reflect.TypeOf(typeToRegister)),
			)

			// registration should not panic
			Register(typeToRegister)

			// new type must be registered
			require.True(
				t,
				TypeIsRegistered(reflect.TypeOf(typeToRegister)),
			)

			// duplicated registration MUST panic
			require.Panics(
				t, func() {
					Register(typeToRegister)
				},
			)
		},
	)

	t.Run(
		"",
		func(t *testing.T) {
			t.Parallel()
			require.Panics(
				t, func() {
					TypeIsRegistered(reflect.TypeOf(0))
				},
			)
		},
	)

	t.Run(
		"",
		func(t *testing.T) {
			t.Parallel()
			require.Panics(
				t, func() {
					Register(0)
				},
			)
		},
	)
}
