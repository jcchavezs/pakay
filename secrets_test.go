package pakay

import (
	"context"
	"log/slog"
	"testing"

	"github.com/jcchavezs/pakay/internal/secrets"
	"github.com/jcchavezs/pakay/internal/sources/env"
	"github.com/stretchr/testify/require"
)

type recordHandler struct {
	records []slog.Record
}

func (rh *recordHandler) Enabled(context.Context, slog.Level) bool { return true }
func (rh *recordHandler) Handle(_ context.Context, r slog.Record) error {
	rh.records = append(rh.records, r)
	return nil
}
func (rh *recordHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return rh }
func (rh *recordHandler) WithGroup(name string) slog.Handler       { return rh }

func unloadSecrets() {
	secrets.All = make(map[string]secrets.Secret)
	secrets.Loaded = false
}

func TestLoadSecretsConfig(t *testing.T) {
	t.Run("secrets are not loaded yet", func(t *testing.T) {
		val, ok := GetSecret(context.Background(), "test_secret")
		require.False(t, ok)
		require.Empty(t, val)

		_, err := AssertSecrets(context.Background())
		require.Error(t, err)
		require.ErrorContains(t, err, "secrets haven't been loaded yet")
	})

	t.Run("loads secrets successfully", func(t *testing.T) {
		t.Cleanup(unloadSecrets)

		RegisterSource(env.Source)

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

		err := LoadSecretsConfig([]byte(config))
		require.NoError(t, err)

		val, ok := GetSecret(context.Background(), "test_secret_1")
		require.True(t, ok)
		require.Equal(t, "test_value", val)

		_, ok = GetSecret(context.Background(), "test_secret_2")
		require.False(t, ok)
	})

	t.Run("unknown secret", func(t *testing.T) {
		t.Cleanup(unloadSecrets)

		config := `---`

		lh := &recordHandler{}

		err := LoadSecretsConfigWithOptions([]byte(config), LoadConfigOptions{
			LoadOptions: LoadOptions{LogHandler: lh},
		})
		require.NoError(t, err)

		_, ok := GetSecret(context.Background(), "unknown_secret")
		require.False(t, ok)
		require.Len(t, lh.records, 1)
		require.Equal(t, "Unknown secret", lh.records[0].Message)
		lh.records[0].Attrs(func(attr slog.Attr) bool {
			require.Equal(t, "name", attr.Key)
			require.Equal(t, "unknown_secret", attr.Value.String())
			return true
		})
	})

	t.Run("renders template variables", func(t *testing.T) {
		t.Cleanup(unloadSecrets)

		RegisterSource(env.Source)

		config := `---
- name: test_secret
  sources:
  - type: env
    env:
      key: {{ $.EnvKey }}
`
		opt := LoadConfigOptions{
			Variables: map[string]string{
				"EnvKey": "TEST_ENV_VAR",
			},
		}

		t.Setenv("TEST_ENV_VAR", "test_value")

		err := LoadSecretsConfigWithOptions([]byte(config), opt)
		require.NoError(t, err)

		val, ok := GetSecret(context.Background(), "test_secret")
		require.True(t, ok)
		require.Equal(t, "test_value", val)
	})

	t.Run("returns error for unknown source", func(t *testing.T) {
		t.Cleanup(unloadSecrets)

		config := `---
- name: test_secret
  sources:
  - type: unknown_source
`

		err := LoadSecretsConfig([]byte(config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown source: unknown_source")
	})

	t.Run("returns error for duplicated secret", func(t *testing.T) {
		t.Cleanup(unloadSecrets)

		config := `---
- name: test_secret
  sources:
  - type: env
    env:
      key: MY_VAR_1
- name: test_secret
  sources:
  - type: env
    env:
      key: MY_VAR_2
`

		err := LoadSecretsConfig([]byte(config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicated declaration for \"test_secret\"")
	})
}
