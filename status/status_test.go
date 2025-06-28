package status

import (
	"context"
	"slices"
	"testing"

	"github.com/jcchavezs/pakay"
	_ "github.com/jcchavezs/pakay/internal/providers"
	"github.com/stretchr/testify/require"
)

func TestCheckSecrets(t *testing.T) {
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
`

	err := pakay.LoadSecretsFromBytes([]byte(config))
	require.NoError(t, err)

	ctx := context.Background()

	ss := CheckSecretsWithOptions(ctx, CheckOptions{
		FilterIn: func(s pakay.Source) bool {
			return !slices.Contains(s.Labels, "deprecated")
		},
	})

	require.Len(t, ss, 2)

	for _, s := range ss {
		require.Len(t, ss[0].Sources(), 1)
		if s.Name() == "test_secret_1" {
			require.Equal(t, "This is a test secret", s.Description())
			require.Equal(t, "env: TEST_ENV_VAR_1", s.Sources()[0])
			t.Setenv("DEPRECATED_TEST_ENV_VAR_1", "A")
			v, ok := s.GetValue(ctx)
			require.Equal(t, "", v)
			require.False(t, ok)

			t.Setenv("TEST_ENV_VAR_1", "A")
			v, ok = s.GetValue(ctx)
			require.Equal(t, "A", v)
			require.True(t, ok)
		} else {
			require.Equal(t, "env: TEST_ENV_VAR_2", s.Sources()[0])
		}
	}
}
