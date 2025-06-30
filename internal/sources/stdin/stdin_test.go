package stdin

import (
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
}
