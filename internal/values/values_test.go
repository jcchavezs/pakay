package values

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFromMap(t *testing.T) {
	t.Run("returns value when key exists and is a string", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		}

		val, ok := GetFromMap[string](m, "key1")
		require.True(t, ok)
		require.Equal(t, "value1", val)
	})

	t.Run("returns false when key does not exist", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": "value1",
		}

		val, ok := GetFromMap[string](m, "key2")
		require.False(t, ok)
		require.Empty(t, val)
	})

	t.Run("returns false when key exists but value is not a string", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": 123,
		}

		val, ok := GetFromMap[string](m, "key1")
		require.False(t, ok)
		require.Empty(t, val)
	})
}
