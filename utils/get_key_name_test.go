package utils

import (
	"reflect"
	"testing"
)

func Test_getKeyName(t *testing.T) {
	type args struct {
		rootKey   string
		fieldType reflect.StructField
	}

	tests := []struct {
		name string
		args args
		want string
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
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := GetKeyName(tt.args.rootKey, tt.args.fieldType); got != tt.want {
					t.Errorf("GetKeyName() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_appendKey(t *testing.T) {
	type args struct {
		a string
		b string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when root",
			args: args{
				a: "",
				b: "xxx",
			},
			want: "xxx",
		}, {
			name: "otherwise",
			args: args{
				a: "xxx",
				b: "yyy",
			},
			want: "xxx-yyy",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := appendKey(tt.args.a, tt.args.b); got != tt.want {
					t.Errorf("appendKey() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
