package pakay

import (
	"testing"

	"github.com/jcchavezs/pakay/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestSecretsConfigToManifestEntriesParity(t *testing.T) {
	manifest := `---
- name: static_secret
  description: Static secret
  sources:
    - type: static
      static:
        value: hello_world
- name: bash_secret
  description: Bash secret
  sources:
    - type: bash
      bash:
        command: echo hi
        timeout_ms: 1000
- name: env_secret
  sources:
    - type: env
      env:
        key: ENV_KEY
- name: stdin_secret
  description: Stdin secret
  sources:
    - type: stdin
      stdin:
        prompt: enter value
- name: op_secret
  sources:
    - type: 1password
      1password:
        ref: op://vault/item/field
`

	programmatic := SecretsConfig{
		{
			Name:        "static_secret",
			Description: "Static secret",
			Sources: []SecretSource{{
				TypedConfig: &StaticConfig{Value: "hello_world"},
			}},
		},
		{
			Name:        "bash_secret",
			Description: "Bash secret",
			Sources: []SecretSource{{
				TypedConfig: &BashConfig{Command: "echo hi", TimeoutMS: 1000},
			}},
		},
		{
			Name: "env_secret",
			Sources: []SecretSource{{
				TypedConfig: &EnvConfig{Key: "ENV_KEY"},
			}},
		},
		{
			Name:        "stdin_secret",
			Description: "Stdin secret",
			Sources: []SecretSource{{
				TypedConfig: &StdinConfig{Prompt: "enter value"},
			}},
		},
		{
			Name: "op_secret",
			Sources: []SecretSource{{
				TypedConfig: &OnePasswordConfig{Ref: "op://vault/item/field"},
			}},
		},
	}.toManifestEntries()

	parsed, err := parser.ParseManifest([]byte(manifest), nil)
	require.NoError(t, err)

	require.Equal(t, parsed, programmatic)
}
