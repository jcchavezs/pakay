package static

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_String(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		config := Config{Value: ""}
		require.Equal(t, "", config.String())
	})

	t.Run("single character", func(t *testing.T) {
		config := Config{Value: "a"}
		require.Equal(t, "*", config.String())
	})

	t.Run("two characters", func(t *testing.T) {
		config := Config{Value: "ab"}
		require.Equal(t, "**", config.String())
	})

	t.Run("three characters", func(t *testing.T) {
		config := Config{Value: "abc"}
		require.Equal(t, "***", config.String())
	})

	t.Run("four characters", func(t *testing.T) {
		config := Config{Value: "abcd"}
		require.Equal(t, "****", config.String())
	})

	t.Run("six characters", func(t *testing.T) {
		config := Config{Value: "abcdef"}
		require.Equal(t, "abc***", config.String())
	})

	t.Run("longer value", func(t *testing.T) {
		config := Config{Value: "this_is_a_long_secret"}
		require.Equal(t, "thi******************", config.String())
	})
}

func TestSource_ConfigFactory(t *testing.T) {
	config := Source.ConfigFactory()
	require.IsType(t, &Config{}, config)
}

func TestSource_SecretGetterFactory(t *testing.T) {
	t.Run("valid config with value", func(t *testing.T) {
		testValue := "test_secret_value"
		config := &Config{
			Value: testValue,
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.True(t, ok)
		require.Equal(t, testValue, secret)
	})

	t.Run("valid config with empty value", func(t *testing.T) {
		config := &Config{
			Value: "",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.True(t, ok)
		require.Empty(t, secret)
	})

	t.Run("invalid config type", func(t *testing.T) {
		invalidConfig := bytes.NewBufferString("invalid")

		getter, err := Source.SecretGetterFactory(invalidConfig)
		require.Error(t, err)
		require.Nil(t, getter)
		require.Equal(t, "invalid config", err.Error())
	})
}
