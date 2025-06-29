package view

import (
	"context"
	"slices"
	"testing"

	"github.com/jcchavezs/pakay"
	_ "github.com/jcchavezs/pakay/internal/sources"
	"github.com/stretchr/testify/require"
)

func TestListSecrets(t *testing.T) {
	config := `---
- name: test_secret_1
  description: This is a test secret
  sources:
    - type: env
      env:
        key: TEST_ENV_VAR_1
    - type: env
      labels: [deprecated]
      env:
        key: DEPRECATED_TEST_ENV_VAR_1
- name: test_secret_2
  sources:
  - type: env
    env:
      key: TEST_ENV_VAR_2
- name: test_secret_3
  sources:
  - type: env
    labels: [deprecated]
    env:
      key: TEST_ENV_VAR_3
`

	err := pakay.LoadSecretsFromBytes([]byte(config))
	require.NoError(t, err)

	ctx := context.Background()

	ss := ListSecretsWithOptions(ctx, GetOptions{
		FilterIn: func(s pakay.Source) bool {
			return !slices.Contains(s.Labels, "deprecated")
		},
	})

	require.Len(t, ss, 3)

	for _, s := range ss {
		switch s.Name() {
		case "test_secret_1":
			require.Len(t, s.Sources(), 1)
			require.Equal(t, "This is a test secret", s.Description())
			require.Equal(t, "env: TEST_ENV_VAR_1", s.Sources()[0])
			t.Setenv("DEPRECATED_TEST_ENV_VAR_1", "my_value")
			v, ok := s.GetValue(ctx)
			require.Equal(t, "", v)
			require.False(t, ok)

			t.Setenv("TEST_ENV_VAR_1", "my_value")
			v, ok = s.GetValue(ctx)
			require.Equal(t, "my_value", v)
			require.True(t, ok)
		case "test_secret_2":
			require.Len(t, s.Sources(), 1)
			require.Equal(t, "env: TEST_ENV_VAR_2", s.Sources()[0])
		case "test_secret_3":
			require.Len(t, s.Sources(), 0)
		}
	}
}
