package env

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_String(t *testing.T) {
	config := Config{
		Key: "TEST_ENV_VAR",
	}

	require.Equal(t, "TEST_ENV_VAR", config.String())
}

func TestProvider_ConfigFactory(t *testing.T) {
	config := Provider.ConfigFactory()
	require.IsType(t, &Config{}, config)
}

func TestProvider_SecretGetterFactory(t *testing.T) {
	t.Run("valid config with existing environment variable", func(t *testing.T) {
		// Set up test environment variable
		testKey := "TEST_EXISTING_VAR"
		testValue := "test_secret_value"
		t.Setenv(testKey, testValue)

		config := &Config{
			Key: testKey,
		}

		getter, err := Provider.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.True(t, ok)
		require.Equal(t, testValue, secret)
	})

	t.Run("valid config with non-existing environment variable", func(t *testing.T) {
		config := &Config{
			Key: "NON_EXISTING_VAR",
		}

		getter, err := Provider.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.False(t, ok)
		require.Empty(t, secret)
	})

	t.Run("valid config with empty environment variable", func(t *testing.T) {
		testKey := "TEST_EMPTY_VAR"
		t.Setenv(testKey, "")

		config := &Config{
			Key: testKey,
		}

		getter, err := Provider.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.False(t, ok)
		require.Empty(t, secret)
	})

	t.Run("valid config with whitespace-only environment variable", func(t *testing.T) {
		testKey := "TEST_WHITESPACE_VAR"
		testValue := "   "
		t.Setenv(testKey, testValue)

		config := &Config{
			Key: testKey,
		}

		getter, err := Provider.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)

		ctx := context.Background()
		secret, ok := getter(ctx)
		require.True(t, ok)
		require.Equal(t, testValue, secret)
	})

	t.Run("invalid config type", func(t *testing.T) {
		invalidConfig := bytes.NewBufferString("invalid")

		getter, err := Provider.SecretGetterFactory(invalidConfig)
		require.Error(t, err)
		require.Nil(t, getter)
		require.Equal(t, "invalid config", err.Error())
	})
}
