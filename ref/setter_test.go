package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestR(t *testing.T) {
	t.Parallel()

	t.Run("string_value", func(t *testing.T) {
		t.Parallel()

		value := "test"
		result := R(value)

		assert.NotNil(t, result)
		assert.Equal(
			t,
			&value,
			result,
		)
		assert.Equal(
			t,
			"test",
			*result,
		)
	})

	t.Run("int_value", func(t *testing.T) {
		t.Parallel()

		value := 42
		result := R(value)

		assert.NotNil(t, result)
		assert.Equal(
			t,
			&value,
			result,
		)
		assert.Equal(
			t,
			42,
			*result,
		)
	})

	t.Run("bool_value", func(t *testing.T) {
		t.Parallel()

		value := true
		result := R(value)

		assert.NotNil(t, result)
		assert.Equal(
			t,
			&value,
			result,
		)
		assert.Equal(
			t,
			true,
			*result,
		)
	})

	t.Run("float_value", func(t *testing.T) {
		t.Parallel()

		value := 3.14
		result := R(value)

		assert.NotNil(t, result)
		assert.Equal(
			t,
			&value,
			result,
		)
		assert.Equal(
			t,
			3.14,
			*result,
		)
	})

	t.Run("zero_values", func(t *testing.T) {
		t.Parallel()

		t.Run("zero_string", func(t *testing.T) {
			t.Parallel()

			value := ""
			result := R(value)

			assert.NotNil(t, result)
			assert.Equal(
				t,
				"",
				*result,
			)
		})

		t.Run("zero_int", func(t *testing.T) {
			t.Parallel()

			value := 0
			result := R(value)

			assert.NotNil(t, result)
			assert.Equal(
				t,
				0,
				*result,
			)
		})

		t.Run("zero_bool", func(t *testing.T) {
			t.Parallel()

			value := false
			result := R(value)

			assert.NotNil(t, result)
			assert.Equal(
				t,
				false,
				*result,
			)
		})
	})

	t.Run("struct_value", func(t *testing.T) {
		t.Parallel()

		type testStruct struct {
			Name string
			Age  int
		}

		value := testStruct{Name: "John", Age: 30}
		result := R(value)

		assert.NotNil(t, result)
		assert.Equal(
			t,
			&value,
			result,
		)
		assert.Equal(
			t,
			"John",
			result.Name,
		)
		assert.Equal(
			t,
			30,
			result.Age,
		)
	})

	t.Run("slice_value", func(t *testing.T) {
		t.Parallel()

		value := []string{"a", "b", "c"}
		result := R(value)

		assert.NotNil(t, result)
		assert.Equal(
			t,
			&value,
			result,
		)
		assert.Equal(
			t,
			[]string{"a", "b", "c"},
			*result,
		)
	})

	t.Run("map_value", func(t *testing.T) {
		t.Parallel()

		value := map[string]int{"a": 1, "b": 2}
		result := R(value)

		assert.NotNil(t, result)
		assert.Equal(
			t,
			&value,
			result,
		)
		assert.Equal(
			t,
			map[string]int{"a": 1, "b": 2},
			*result,
		)
	})
}
