package svalue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("create_value", func(t *testing.T) {
		t.Parallel()

		v := Value{
			Location: "env:TEST_VAR",
			Value:    "test_value",
		}

		assert.Equal(
			t,
			"env:TEST_VAR",
			v.Location,
		)
		assert.Equal(
			t,
			"test_value",
			v.Value,
		)
	})

	t.Run("empty_value", func(t *testing.T) {
		t.Parallel()

		v := Value{}

		assert.Equal(
			t,
			"",
			v.Location,
		)
		assert.Equal(
			t,
			"",
			v.Value,
		)
	})

	t.Run("value_with_special_characters", func(t *testing.T) {
		t.Parallel()

		v := Value{
			Location: "file:/path/to/config.yaml:line:15",
			Value:    "value with spaces and symbols !@#$%",
		}

		assert.Equal(
			t,
			"file:/path/to/config.yaml:line:15",
			v.Location,
		)
		assert.Equal(
			t,
			"value with spaces and symbols !@#$%",
			v.Value,
		)
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

	t.Run("add_single_value", func(t *testing.T) {
		t.Parallel()

		values := Values{
			"key1": &Value{
				Location: "env:KEY1",
				Value:    "value1",
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
			"key1",
		)
		assert.Equal(
			t,
			"env:KEY1",
			values["key1"].Location,
		)
		assert.Equal(
			t,
			"value1",
			values["key1"].Value,
		)
	})

	t.Run("add_multiple_values", func(t *testing.T) {
		t.Parallel()

		values := Values{
			"key1": &Value{
				Location: "env:KEY1",
				Value:    "value1",
			},
			"key2": &Value{
				Location: "file:config.yaml",
				Value:    "value2",
			},
			"key3": &Value{
				Location: "cmdline:--key3",
				Value:    "value3",
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
			"key1",
		)
		assert.Contains(
			t,
			values,
			"key2",
		)
		assert.Contains(
			t,
			values,
			"key3",
		)

		// Check values.
		assert.Equal(
			t,
			"value1",
			values["key1"].Value,
		)
		assert.Equal(
			t,
			"value2",
			values["key2"].Value,
		)
		assert.Equal(
			t,
			"value3",
			values["key3"].Value,
		)

		// Check locations.
		assert.Equal(
			t,
			"env:KEY1",
			values["key1"].Location,
		)
		assert.Equal(
			t,
			"file:config.yaml",
			values["key2"].Location,
		)
		assert.Equal(
			t,
			"cmdline:--key3",
			values["key3"].Location,
		)
	})

	t.Run("nil_value_pointer", func(t *testing.T) {
		t.Parallel()

		values := Values{
			"key1": nil,
		}

		assert.Equal(
			t,
			1,
			len(values),
		)
		assert.Contains(
			t,
			values,
			"key1",
		)
		assert.Nil(t, values["key1"])
	})

	t.Run("overwrite_existing_key", func(t *testing.T) {
		t.Parallel()

		values := Values{
			"key1": &Value{
				Location: "env:KEY1",
				Value:    "original_value",
			},
		}

		// Overwrite with new value.
		values["key1"] = &Value{
			Location: "file:config.yaml",
			Value:    "new_value",
		}

		assert.Equal(
			t,
			1,
			len(values),
		)
		assert.Equal(
			t,
			"new_value",
			values["key1"].Value,
		)
		assert.Equal(
			t,
			"file:config.yaml",
			values["key1"].Location,
		)
	})

	t.Run("access_nonexistent_key", func(t *testing.T) {
		t.Parallel()

		values := Values{}

		assert.Nil(t, values["nonexistent"])
	})

	t.Run("complex_key_names", func(t *testing.T) {
		t.Parallel()

		values := Values{
			"database.host": &Value{
				Location: "env",
				Value:    "localhost",
			},
			"database.port":       &Value{Location: "env", Value: "5432"},
			"app.feature.enabled": &Value{Location: "file", Value: "true"},
			"logging.level": &Value{
				Location: "cmdline",
				Value:    "debug",
			},
			"special-chars_123.test": &Value{Location: "test", Value: "value"},
		}

		assert.Equal(
			t,
			5,
			len(values),
		)
		assert.Equal(
			t,
			"localhost",
			values["database.host"].Value,
		)
		assert.Equal(
			t,
			"5432",
			values["database.port"].Value,
		)
		assert.Equal(
			t,
			"true",
			values["app.feature.enabled"].Value,
		)
		assert.Equal(
			t,
			"debug",
			values["logging.level"].Value,
		)
		assert.Equal(
			t,
			"value",
			values["special-chars_123.test"].Value,
		)
	})
}
