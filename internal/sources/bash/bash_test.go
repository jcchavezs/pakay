package bash

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_String(t *testing.T) {
	config := Config{
		Command: "echo 'test'",
	}

	require.Equal(t, "echo 'test'", config.String())
}

func TestSource_ConfigFactory(t *testing.T) {
	config := Source.ConfigFactory()
	require.IsType(t, &Config{}, config)
}

func TestSource_SecretGetterFactory(t *testing.T) {
	t.Run("valid config with successful command", func(t *testing.T) {
		config := &Config{
			Command: `echo "test_secret"`,
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.True(t, ok)
		require.Equal(t, "test_secret", secret)
	})

	t.Run("valid config with failing command", func(t *testing.T) {
		config := &Config{
			Command: "exit 1",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.False(t, ok)
		require.Empty(t, secret)
	})

	t.Run("valid config with timeout", func(t *testing.T) {
		config := &Config{
			Command:   "sleep 2 && echo 'delayed_secret'",
			TimeoutMS: 100,
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.False(t, ok)
		require.Empty(t, secret)
	})

	t.Run("trims whitespace from output", func(t *testing.T) {
		config := &Config{
			Command: `echo "  test_with_spaces  "`,
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.True(t, ok)
		require.Equal(t, "test_with_spaces", secret)
	})

	t.Run("empty ref returns error", func(t *testing.T) {
		config := &Config{
			Command: "",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.Error(t, err)
		require.Nil(t, getter)
		require.Equal(t, "command cannot be empty", err.Error())
	})

	t.Run("invalid config type", func(t *testing.T) {
		invalidConfig := bytes.NewBufferString("invalid")

		getter, err := Source.SecretGetterFactory(invalidConfig)
		require.Error(t, err)
		require.Nil(t, getter)
		require.Equal(t, "invalid config", err.Error())
	})
}
