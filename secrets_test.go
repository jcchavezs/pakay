package pakay

import (
	"context"
	"log/slog"
	"slices"
	"testing"

	"github.com/jcchavezs/pakay/internal/parser"
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

		err := ParseAndLoadSecrets([]byte(config))
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

		err := ParseAndLoadSecretsWithOptions([]byte(config), ParseAndLoadOptions{
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
		opt := ParseAndLoadOptions{
			Variables: map[string]string{
				"EnvKey": "TEST_ENV_VAR",
			},
		}

		t.Setenv("TEST_ENV_VAR", "test_value")

		err := ParseAndLoadSecretsWithOptions([]byte(config), opt)
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

		err := ParseAndLoadSecrets([]byte(config))
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

		err := ParseAndLoadSecrets([]byte(config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicated declaration for \"test_secret\"")
	})
}

func TestAssertSecretsWithOptions(t *testing.T) {
	// Save original state
	originalAll := secrets.All
	originalLoaded := secrets.Loaded

	// Restore original state after test
	t.Cleanup(func() {
		secrets.All = originalAll
		secrets.Loaded = originalLoaded
	})

	t.Run("returns error when secrets not loaded", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{}
		secrets.Loaded = false

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{})

		require.Error(t, err)
		require.Nil(t, missing)
		require.Equal(t, "secrets haven't been loaded yet", err.Error())
	})

	t.Run("returns empty list when all secrets are available", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{
			"secret1": {
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "dev"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "value1", true
						},
					},
				},
			},
			"secret2": {
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "prod"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "value2", true
						},
					},
				},
			},
		}
		secrets.Loaded = true

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{})

		require.NoError(t, err)
		require.Empty(t, missing)
	})

	t.Run("returns list of missing secrets", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{
			"available_secret": {
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "dev"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "value", true
						},
					},
				},
			},
			"missing_secret": {
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "prod"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "", false
						},
					},
				},
			},
		}
		secrets.Loaded = true

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{})

		require.NoError(t, err)
		require.Len(t, missing, 1)
		require.Contains(t, missing, "missing_secret")
	})

	t.Run("filters secrets by SecretFilterIn", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{
			"dev_secret": {
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "dev"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "", false
						},
					},
				},
			},
			"prod_secret": {
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "prod"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "", false
						},
					},
				},
			},
		}
		secrets.Loaded = true

		// Filter to only check secrets with "dev" in the name
		secretFilter := func(s Secret) bool {
			return s.Name == "dev_secret"
		}

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{
			SecretFilterIn: secretFilter,
		})

		require.NoError(t, err)
		require.Len(t, missing, 1)
		require.Contains(t, missing, "dev_secret")
		require.NotContains(t, missing, "prod_secret")
	})

	t.Run("filters sources by SourceFilterIn", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{
			"multi_source_secret": {
				ManifestEntry: parser.ManifestEntry{
					Name: "multi_source_secret",
					Sources: []parser.ManifestEntrySource{
						{Type: "env", Labels: []string{"env", "dev"}},
						{Type: "static", Labels: []string{"env", "prod"}},
					},
				},
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "dev"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "", false // dev source fails
						},
					},
					{
						Labels: []string{"env", "prod"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "prod_value", true // prod source succeeds
						},
					},
				},
			},
		}
		secrets.Loaded = true

		// Filter to only use prod sources
		sourceFilter := func(s Source) bool {
			for _, label := range s.Labels {
				if label == "prod" {
					return true
				}
			}
			return false
		}

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{
			SourceFilterIn: sourceFilter,
		})

		require.NoError(t, err)
		require.Empty(t, missing) // Should be empty because prod source succeeds
	})

	t.Run("combines SecretFilterIn and SourceFilterIn", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{
			"dev_secret": {
				ManifestEntry: parser.ManifestEntry{
					Name: "dev_secret",
					Sources: []parser.ManifestEntrySource{
						{Type: "env", Labels: []string{"env", "dev"}},
					},
				},
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "dev"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "dev_value", true
						},
					},
				},
			},
			"prod_secret": {
				ManifestEntry: parser.ManifestEntry{
					Name: "prod_secret",
					Sources: []parser.ManifestEntrySource{
						{Type: "env", Labels: []string{"env", "prod"}},
					},
				},
				Getters: []secrets.Getter{
					{
						Labels: []string{"env", "prod"},
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "", false
						},
					},
				},
			},
		}
		secrets.Loaded = true

		// Only check dev secrets
		secretFilter := func(s Secret) bool {
			return s.Name == "dev_secret"
		}

		// Only use dev sources
		sourceFilter := func(s Source) bool {
			return slices.Contains(s.Labels, "dev")
		}

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{
			SecretFilterIn: secretFilter,
			SourceFilterIn: sourceFilter,
		})

		require.NoError(t, err)
		require.Empty(t, missing) // dev_secret with dev source is available
	})

	t.Run("empty secrets map", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{}
		secrets.Loaded = true

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{})

		require.NoError(t, err)
		require.Empty(t, missing)
	})

	t.Run("nil filter options", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{
			"test_secret": {
				Getters: []secrets.Getter{
					{
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "value", true
						},
					},
				},
			},
		}
		secrets.Loaded = true

		ctx := context.Background()
		missing, err := AssertSecretsWithOptions(ctx, AssertOptions{
			SecretFilterIn: nil,
			SourceFilterIn: nil,
		})

		require.NoError(t, err)
		require.Empty(t, missing)
	})
}

func TestAssertSecrets(t *testing.T) {
	// Save original state
	originalAll := secrets.All
	originalLoaded := secrets.Loaded

	// Restore original state after test
	t.Cleanup(func() {
		secrets.All = originalAll
		secrets.Loaded = originalLoaded
	})

	t.Run("calls AssertSecretsWithOptions with empty options", func(t *testing.T) {
		secrets.All = map[string]secrets.Secret{
			"test_secret": {
				Getters: []secrets.Getter{
					{
						SecretGetter: func(ctx context.Context) (string, bool) {
							return "value", true
						},
					},
				},
			},
		}
		secrets.Loaded = true

		ctx := context.Background()
		missing, err := AssertSecrets(ctx)

		require.NoError(t, err)
		require.Empty(t, missing)
	})
}
