package sbased2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errFailApply = errors.New("")

type validOption struct{}

type failOption struct{}

func (validOption) apply(*internalOpts) error {
	return nil
}

func TestWithAliases(t *testing.T) {
	t.Parallel()

	type args struct {
		mapping map[string]string
	}

	tests := []struct {
		args args
		want AliasesOption
		name string
	}{
		{
			name: "empty map returns nil",
			args: args{
				mapping: make(map[string]string),
			},
			want: nil,
		},
		{
			name: "success",
			args: args{
				mapping: map[string]string{"a": "b"},
			},
			want: AliasesOption{"a": "b"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()
				got := WithAliases(tt.args.mapping)
				require.Equalf(
					t, tt.want, got, "WithAliases() = %v, want %v", got,
					tt.want,
				)
			},
		)
	}
}

func TestAliasesOption_apply(t *testing.T) {
	t.Parallel()

	t.Run(
		"nil alias",
		func(t *testing.T) {
			t.Parallel()

			var ao AliasesOption

			o := &internalOpts{}

			require.NoError(t, ao.apply(o))
			require.Nil(t, o.aliases)
		},
	)

	t.Run(
		"success",
		func(t *testing.T) {
			t.Parallel()

			ao := AliasesOption{
				"a": "b",
				"c": "d",
			}

			options := &internalOpts{}

			require.NoError(t, ao.apply(options))
			require.Equal(
				t, map[string]string{
					"a": "b",
					"c": "d",
				}, options.aliases,
			)
		},
	)
}

func (failOption) apply(*internalOpts) error {
	return errFailApply
}

func Test_internalOpts_applyOptions(t *testing.T) {
	t.Parallel()

	type fields struct {
		aliases map[string]string
	}

	type args struct {
		os []Option
	}

	tests := []struct {
		wantErr error
		fields  fields
		name    string
		args    args
	}{
		{
			name: "nil option list",
			fields: fields{
				aliases: nil,
			},
			args: args{
				os: nil,
			},
			wantErr: nil,
		},
		{
			name: "empty options list",
			fields: fields{
				aliases: nil,
			},
			args: args{
				os: nil,
			},
			wantErr: nil,
		},
		{
			name: "no errors",
			fields: fields{
				aliases: nil,
			},
			args: args{
				os: []Option{validOption{}, validOption{}, validOption{}},
			},
			wantErr: nil,
		},
		{
			name: "error in starting position",
			fields: fields{
				aliases: nil,
			},
			args: args{
				os: []Option{failOption{}, validOption{}, validOption{}},
			},
			wantErr: errFailApply,
		},
		{
			name: "error in middle",
			fields: fields{
				aliases: nil,
			},
			args: args{
				os: []Option{
					validOption{}, validOption{}, failOption{}, validOption{},
				},
			},
			wantErr: errFailApply,
		},
		{
			name: "error at last position",
			fields: fields{
				aliases: nil,
			},
			args: args{
				os: []Option{
					validOption{}, validOption{}, failOption{}, validOption{},
				},
			},
			wantErr: errFailApply,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				o := &internalOpts{
					aliases: tt.fields.aliases,
				}
				err := o.applyOptions(tt.args.os)
				require.ErrorIsf(
					t,
					err,
					tt.wantErr,
					"applyOptions() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			},
		)
	}
}
