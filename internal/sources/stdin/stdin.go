package stdin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jcchavezs/pakay/internal/log"

	"github.com/jcchavezs/pakay/types"
	"golang.org/x/term"
)

type Config struct {
	Prompt string `yaml:"prompt"`
}

func (c Config) String() string {
	return "prompt"
}

var readPassword func(int) ([]byte, error) = term.ReadPassword

var Source = types.SecretSource{
	ConfigFactory: func() types.SourceConfig {
		return &Config{}
	},
	SecretGetterFactory: func(cfg types.SourceConfig) (types.SecretGetter, error) {
		var prompt string
		if tCfg, ok := cfg.(*Config); ok {
			prompt = tCfg.Prompt
		} else {
			return nil, errors.New("invalid config")
		}

		if prompt == "" {
			return nil, errors.New("prompt cannot be empty")
		}

		return func(ctx context.Context) (string, bool) {
			_, _ = fmt.Printf("%s: ", prompt)
			input, err := readPassword(int(os.Stdin.Fd()))
			_, _ = fmt.Println("")
			if err != nil {
				log.Logger.Error("failed to read from stdin", "error", err)
				return "", false
			}

			input = bytes.TrimSpace(input)
			return string(input), len(input) > 0
		}, nil
	},
}
