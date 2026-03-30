package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_reflectPlayGround1(t *testing.T) { //nolint:paralleltest
	var pNil *int

	v := 5
	pOK := &v

	vOK := reflect.ValueOf(pOK)
	vNil := reflect.ValueOf(pNil)

	require.Equal(
		t,
		reflect.Ptr,
		vOK.Type().Kind(),
	)
	require.Equal(
		t,
		reflect.Ptr,
		vNil.Type().Kind(),
	)

	require.False(t, vOK.IsNil())
	require.True(t, vNil.IsNil())
}

func deref[T any](v T) *T {
	return &v
}

func Test_reflectPlayGround2(t *testing.T) { //nolint:paralleltest
	// v := []int{1,2,3,4,5}
	// v := []*int{R(1), nil, R(1)}

	type Embedded struct {
		E *int
	}

	type Other struct {
		Embedded
		X *int `json:"x" yaml:"x"`
		Y *int `json:"y" yaml:"y"`
		L []*int
	}

	v := map[string]*Other{
		"A": {
			Embedded: Embedded{
				E: deref(101),
			},
			X: deref(1),
			Y: deref(2),
			L: []*int{
				deref(1),
				deref(2),
			},
		},
		"B": {
			X: deref(1),
			Y: deref(2),
			L: []*int{
				deref(1),
				deref(2),
			},
		},
		"C": {
			Embedded: Embedded{},
			X:        deref(1),
			Y:        deref(2),
			L: []*int{
				deref(1),
				deref(2),
			},
		},
		"D": {
			Embedded: Embedded{
				E: deref(101),
			},
			Y: deref(2),
			L: []*int{
				deref(1),
				deref(2),
			},
		},
		"E": {
			Embedded: Embedded{
				E: deref(101),
			},
			X: deref(1),
			Y: deref(2),
			L: []*int{
				nil,
				deref(2),
			},
		},
	}

	// var v *int

	perfect("", reflect.ValueOf(v))
}

func perfect(path string, value reflect.Value) {
	// fmt.Println(">>", path)

	ov := value
	v := value
	// check for nillness
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			fmt.Println(path, "is nil")
		}

		// follow reference
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			perfect(fmt.Sprintf(
				"%s[%d]",
				path,
				i,
			), v.Index(i))
		}

		return

	case reflect.Map:
		for _, key := range v.MapKeys() {
			perfect(fmt.Sprintf(
				"%s{%v}",
				path,
				key.Interface(),
			), v.MapIndex(key))
		}
		return

	case reflect.Struct:
		list, _ := getVisibleFieldList(path, ov.Type())
		for _, embedded := range list {
			vv := v

			for _, index := range embedded.index {
				vv = vv.Field(index)
			}

			perfect(
				fmt.Sprintf(
					"%s.%s",
					path,
					embedded.field.Name,
				),
				vv,
			)
		}

		return

	default:
		return
	}
}
