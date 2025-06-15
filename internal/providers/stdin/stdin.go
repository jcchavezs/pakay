package stdin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jcchavezs/pakay/internal/values"
	"github.com/jcchavezs/pakay/types"
	"golang.org/x/term"
)

var Provider types.SecretProvider = func(cfg types.ProviderConfig) (types.SecretGetter, error) {
	var prompt, found = values.GetFromMap[string](cfg, "prompt")
	if !found {
		return nil, errors.New("missing stdin.prompt value")
	}

	return func(ctx context.Context) (string, bool) {
		_, _ = fmt.Print(prompt + ": ")
		input, err := term.ReadPassword(int(os.Stdin.Fd()))
		_, _ = fmt.Println("")
		if err != nil {
			slog.Error("failed to read from stdin", "error", err)
			return "", false
		}

		return string(input), true
	}, nil
}
