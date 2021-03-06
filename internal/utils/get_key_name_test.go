package utils

import (
	"reflect"
	"testing"
)

func Test_getKeyName(t *testing.T) {
	t.Parallel()

	type args struct {
		rootKey   string
		fieldType reflect.StructField
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "not yaml",
			args: args{
				rootKey: "root",
				fieldType: reflect.StructField{
					Name: `Bolos`,
					Tag:  `json:"toto"`,
				},
			},
			want: "root-bolos",
		},
		{
			name: "none",
			args: args{
				rootKey: "root",
				fieldType: reflect.StructField{
					Name: "Folo",
					Tag:  "",
				},
			},
			want: "root-folo",
		},
		{
			name: "none",
			args: args{
				rootKey: "root",
				fieldType: reflect.StructField{
					Name: "FoloPolo",
					Tag:  "",
				},
			},
			want: "root-folo_polo",
		},

		{
			name: "success",
			args: args{
				rootKey: "root",
				fieldType: reflect.StructField{
					Tag: `yaml:"toto"`,
				},
			},
			want: "root-toto",
		},
		{
			name: "success with no prefix",
			args: args{
				rootKey: "",
				fieldType: reflect.StructField{
					Tag: `yaml:"toto"`,
				},
			},
			want: "toto",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				if got := GetKeyName(
					tt.args.rootKey, tt.args.fieldType,
				); got != tt.want {
					t.Errorf("GetKeyName() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
