package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_String(t *testing.T) {
	config := Config{
		Ref: "op://vault/item/field",
	}

	require.Equal(t, "op://vault/item/field", config.String())
}

func TestSource_SecretGetterFactory(t *testing.T) {
	t.Run("empty ref returns error", func(t *testing.T) {
		config := &Config{
			Ref: "",
		}

		getter, err := Source.SecretGetterFactory(config)
		require.Error(t, err)
		require.Nil(t, getter)
		require.Equal(t, "ref cannot be empty", err.Error())
	})

	t.Run("invalid config type", func(t *testing.T) {
		invalidConfig := bytes.NewBufferString("invalid")

		getter, err := Source.SecretGetterFactory(invalidConfig)
		require.Error(t, err)
		require.Nil(t, getter)
		require.Equal(t, "invalid config", err.Error())
	})
}
