package pakay

import (
	"context"
	"testing"

	"github.com/jcchavezs/pakay/internal/providers/env"
	"github.com/stretchr/testify/require"
)

func TestLoadSecretsFromBytes(t *testing.T) {
	t.Run("loads secrets successfully", func(t *testing.T) {
		RegisterProvider("env", env.Provider)

		config := `---
- name: test_secret_1
  sources:
    - type: env
      env:
        key: TEST_ENV_VAR_1
- name: test_secret_2
  sources:
  - type: env
    env:
      key: TEST_ENV_VAR_2
`

		t.Setenv("TEST_ENV_VAR_1", "test_value")

		err := LoadSecretsFromBytes([]byte(config))
		require.NoError(t, err)

		val, ok := GetSecret(context.Background(), "test_secret_1")
		require.True(t, ok)
		require.Equal(t, "test_value", val)

		_, ok = GetSecret(context.Background(), "test_secret_2")
		require.False(t, ok)
	})

	t.Run("renders template variables", func(t *testing.T) {
		RegisterProvider("env", env.Provider)

		config := `---
- name: test_secret
  sources:
  - type: env
    env:
      key: {{ $.EnvKey }}
`
		opt := LoadOptions{
			Variables: map[string]string{
				"EnvKey": "TEST_ENV_VAR",
			},
		}

		t.Setenv("TEST_ENV_VAR", "test_value")

		err := LoadSecretsFromBytesWithOptions([]byte(config), opt)
		require.NoError(t, err)

		val, ok := GetSecret(context.Background(), "test_secret")
		require.True(t, ok)
		require.Equal(t, "test_value", val)
	})

	t.Run("returns error for unknown provider", func(t *testing.T) {
		config := `---
- name: test_secret
  sources:
  - type: unknown_provider
`

		err := LoadSecretsFromBytes([]byte(config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown provider: unknown_provider")
	})
}
