package stdin

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jcchavezs/pakay/internal/log"

	internaltypes "github.com/jcchavezs/pakay/internal/types"

	"github.com/jcchavezs/pakay/types"
	"golang.org/x/term"
)

type Config struct {
	internaltypes.TypedConfig
	Prompt string `yaml:"prompt"`
}

func (*Config) String() string {
	return "prompt"
}

func (*Config) Type() string {
	return "stdin"
}

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
			_, _ = fmt.Print(prompt + ": ")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			_, _ = fmt.Println("")
			if err != nil {
				log.Logger.Error("failed to read from stdin", "error", err)
				return "", false
			}

			return string(input), true
		}, nil
	},
}
