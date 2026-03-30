package fvalue

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("create_value", func(t *testing.T) {
		t.Parallel()

		stringVal := reflect.ValueOf("test")

		v := Value{
			Value:    stringVal,
			Location: "env:TEST_VAR",
			Path:     "config.field",
		}

		assert.Equal(
			t,
			stringVal,
			v.Value,
		)
		assert.Equal(
			t,
			"env:TEST_VAR",
			v.Location,
		)
		assert.Equal(
			t,
			"config.field",
			v.Path,
		)
	})

	t.Run("empty_value", func(t *testing.T) {
		t.Parallel()

		v := Value{}

		assert.Equal(
			t,
			reflect.Value{},
			v.Value,
		)
		assert.Equal(
			t,
			"",
			v.Location,
		)
		assert.Equal(
			t,
			"",
			v.Path,
		)
	})

	t.Run("different_reflect_types", func(t *testing.T) {
		t.Parallel()

		t.Run("int_value", func(t *testing.T) {
			t.Parallel()

			intVal := reflect.ValueOf(42)
			v := Value{
				Value:    intVal,
				Location: "cmdline:--port",
				Path:     "server.port",
			}

			assert.Equal(
				t,
				intVal,
				v.Value,
			)
			assert.Equal(
				t,
				"cmdline:--port",
				v.Location,
			)
			assert.Equal(
				t,
				"server.port",
				v.Path,
			)
		})

		t.Run("bool_value", func(t *testing.T) {
			t.Parallel()

			boolVal := reflect.ValueOf(true)
			v := Value{
				Value:    boolVal,
				Location: "file:config.yaml:line:10",
				Path:     "app.debug",
			}

			assert.Equal(
				t,
				boolVal,
				v.Value,
			)
			assert.Equal(
				t,
				"file:config.yaml:line:10",
				v.Location,
			)
			assert.Equal(
				t,
				"app.debug",
				v.Path,
			)
		})

		t.Run("slice_value", func(t *testing.T) {
			t.Parallel()

			sliceVal := reflect.ValueOf([]string{"a", "b", "c"})
			v := Value{
				Value:    sliceVal,
				Location: "env:LIST_VAR",
				Path:     "app.items",
			}

			assert.Equal(
				t,
				sliceVal,
				v.Value,
			)
			assert.Equal(
				t,
				"env:LIST_VAR",
				v.Location,
			)
			assert.Equal(
				t,
				"app.items",
				v.Path,
			)
		})

		t.Run("struct_value", func(t *testing.T) {
			t.Parallel()

			type testStruct struct {
				Name string
				Age  int
			}

			structVal := reflect.ValueOf(testStruct{Name: "John", Age: 30})
			v := Value{
				Value:    structVal,
				Location: "struct:input",
				Path:     "user",
			}

			assert.Equal(
				t,
				structVal,
				v.Value,
			)
			assert.Equal(
				t,
				"struct:input",
				v.Location,
			)
			assert.Equal(
				t,
				"user",
				v.Path,
			)
		})
	})

	t.Run("complex_paths", func(t *testing.T) {
		t.Parallel()

		stringVal := reflect.ValueOf("value")

		testCases := []struct {
			name     string
			path     string
			location string
		}{
			{
				name:     "simple_path",
				path:     "field",
				location: "env:FIELD",
			},
			{
				name:     "nested_path",
				path:     "parent.child.grandchild",
				location: "file:config.yaml",
			},
			{
				name:     "array_index_path",
				path:     "items[0].name",
				location: "cmdline:--items",
			},
			{
				name:     "complex_path",
				path:     "database.connections[primary].host",
				location: "env:DB_PRIMARY_HOST",
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				v := Value{
					Value:    stringVal,
					Location: tc.location,
					Path:     tc.path,
				}

				assert.Equal(
					t,
					stringVal,
					v.Value,
				)
				assert.Equal(
					t,
					tc.location,
					v.Location,
				)
				assert.Equal(
					t,
					tc.path,
					v.Path,
				)
			})
		}
	})
}

func TestValues(t *testing.T) {
	t.Parallel()

	t.Run("empty_values", func(t *testing.T) {
		t.Parallel()

		values := Values{}

		assert.NotNil(t, values)
		assert.Equal(
			t,
			0,
			len(values),
		)
	})

	t.Run("single_value", func(t *testing.T) {
		t.Parallel()

		values := Values{
			1: &Value{
				Value:    reflect.ValueOf("test"),
				Location: "env:TEST",
				Path:     "field1",
			},
		}

		assert.Equal(
			t,
			1,
			len(values),
		)
		assert.Contains(
			t,
			values,
			uint(1),
		)
		assert.Equal(
			t,
			"test",
			values[1].Value.String(),
		)
		assert.Equal(
			t,
			"env:TEST",
			values[1].Location,
		)
		assert.Equal(
			t,
			"field1",
			values[1].Path,
		)
	})

	t.Run("multiple_values", func(t *testing.T) {
		t.Parallel()

		values := Values{
			1: &Value{
				Value:    reflect.ValueOf("string_value"),
				Location: "env:STRING_VAR",
				Path:     "config.string_field",
			},
			2: &Value{
				Value:    reflect.ValueOf(42),
				Location: "cmdline:--number",
				Path:     "config.number_field",
			},
			3: &Value{
				Value:    reflect.ValueOf(true),
				Location: "file:config.yaml",
				Path:     "config.bool_field",
			},
		}

		assert.Equal(
			t,
			3,
			len(values),
		)

		// Check all keys exist.
		assert.Contains(
			t,
			values,
			uint(1),
		)
		assert.Contains(
			t,
			values,
			uint(2),
		)
		assert.Contains(
			t,
			values,
			uint(3),
		)

		// Check string value.
		assert.Equal(
			t,
			"string_value",
			values[1].Value.String(),
		)
		assert.Equal(
			t,
			"env:STRING_VAR",
			values[1].Location,
		)
		assert.Equal(
			t,
			"config.string_field",
			values[1].Path,
		)

		// Check int value.
		assert.Equal(
			t,
			int64(42),
			values[2].Value.Int(),
		)
		assert.Equal(
			t,
			"cmdline:--number",
			values[2].Location,
		)
		assert.Equal(
			t,
			"config.number_field",
			values[2].Path,
		)

		// Check bool value.
		assert.Equal(
			t,
			true,
			values[3].Value.Bool(),
		)
		assert.Equal(
			t,
			"file:config.yaml",
			values[3].Location,
		)
		assert.Equal(
			t,
			"config.bool_field",
			values[3].Path,
		)
	})

	t.Run("nil_value_pointer", func(t *testing.T) {
		t.Parallel()

		values := Values{
			1: nil,
		}

		assert.Equal(
			t,
			1,
			len(values),
		)
		assert.Contains(
			t,
			values,
			uint(1),
		)
		assert.Nil(t, values[1])
	})

	t.Run("overwrite_existing_key", func(t *testing.T) {
		t.Parallel()

		values := Values{
			1: &Value{
				Value:    reflect.ValueOf("original"),
				Location: "env:ORIG",
				Path:     "field",
			},
		}

		// Overwrite with new value.
		values[1] = &Value{
			Value:    reflect.ValueOf("new"),
			Location: "file:config",
			Path:     "field",
		}

		assert.Equal(
			t,
			1,
			len(values),
		)
		assert.Equal(
			t,
			"new",
			values[1].Value.String(),
		)
		assert.Equal(
			t,
			"file:config",
			values[1].Location,
		)
		assert.Equal(
			t,
			"field",
			values[1].Path,
		)
	})

	t.Run("access_nonexistent_key", func(t *testing.T) {
		t.Parallel()

		values := Values{}

		assert.Nil(t, values[999])
	})

	t.Run("large_key_values", func(t *testing.T) {
		t.Parallel()

		values := Values{
			1000: &Value{
				Value:    reflect.ValueOf("value1000"),
				Location: "test:1000",
				Path:     "field1000",
			},
			9999: &Value{
				Value:    reflect.ValueOf("value9999"),
				Location: "test:9999",
				Path:     "field9999",
			},
		}

		assert.Equal(
			t,
			2,
			len(values),
		)
		assert.Equal(
			t,
			"value1000",
			values[1000].Value.String(),
		)
		assert.Equal(
			t,
			"value9999",
			values[9999].Value.String(),
		)
	})

	t.Run("zero_key", func(t *testing.T) {
		t.Parallel()

		values := Values{
			0: &Value{
				Value:    reflect.ValueOf("zero_value"),
				Location: "env:ZERO",
				Path:     "zero_field",
			},
		}

		assert.Equal(
			t,
			1,
			len(values),
		)
		assert.Contains(
			t,
			values,
			uint(0),
		)
		assert.Equal(
			t,
			"zero_value",
			values[0].Value.String(),
		)
	})
}
