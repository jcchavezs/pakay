package stdin

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSource_SecretGetterFactory(t *testing.T) {
	t.Run("empty ref returns error", func(t *testing.T) {
		config := &Config{
			Prompt: "",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.Error(t, err)
		require.Nil(t, getter)
		require.Equal(t, "prompt cannot be empty", err.Error())
	})

	t.Run("empty input returns false", func(t *testing.T) {
		readPassword = func(int) ([]byte, error) {
			return nil, nil
		}
		config := &Config{
			Prompt: "Insert the password",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)
		v, ok := getter(context.Background())
		require.Empty(t, v)
		require.False(t, ok)
	})

	t.Run("space only input returns false", func(t *testing.T) {
		readPassword = func(int) ([]byte, error) {
			return []byte{' '}, nil
		}
		config := &Config{
			Prompt: "Insert the password",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)
		v, ok := getter(context.Background())
		require.Empty(t, v)
		require.False(t, ok)
	})

	t.Run("valid input returns true", func(t *testing.T) {
		readPassword = func(int) ([]byte, error) {
			return []byte("my_password"), nil
		}
		config := &Config{
			Prompt: "Insert the password",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)
		v, ok := getter(context.Background())
		require.Equal(t, "my_password", v)
		require.True(t, ok)
	})

	t.Run("invalid input returns true", func(t *testing.T) {
		readPassword = func(int) ([]byte, error) {
			return nil, errors.New("invalid input")
		}
		config := &Config{
			Prompt: "Insert the password",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.NoError(t, err)
		require.NotNil(t, getter)
		v, ok := getter(context.Background())
		require.Empty(t, v)
		require.False(t, ok)
	})
}
