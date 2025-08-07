package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	stdexec "os/exec"
	"strings"

	"github.com/jcchavezs/pakay/internal/exec"
	"github.com/jcchavezs/pakay/internal/log"
	internaltypes "github.com/jcchavezs/pakay/internal/types"
	"github.com/jcchavezs/pakay/types"
)

type Config struct {
	Ref string `yaml:"ref"`
}

func (c *Config) String() string {
	return c.Ref
}

func (*Config) Type() string {
	return "1password"
}

func (*Config) SentinelFn(internaltypes.SentinelVal) {}

var Source = types.SecretSource{
	ConfigFactory: func() types.SourceConfig {
		return &Config{}
	},
	SecretGetterFactory: func(cfg types.SourceConfig) (types.SecretGetter, error) {
		var ref string
		if tCfg, ok := cfg.(*Config); ok {
			ref = tCfg.Ref
		} else {
			return nil, errors.New("invalid config")
		}

		if ref == "" {
			return nil, errors.New("ref cannot be empty")
		}

		return func(ctx context.Context) (string, bool) {
			_, err := stdexec.LookPath("op")
			if err != nil {
				log.Logger.Error("1Password CLI not found", "error", err)
				return "", false
			}

			if out, err := exec.CommandContext(ctx, "op", "account", "list"); err != nil {
				return "", false
			} else if len(bytes.TrimSpace(out)) == 0 {
				_, _ = fmt.Fprintf(os.Stderr, "You can use 1Password by turning on the 1Password desktop app integration by following this instructions:\nhttps://developer.1password.com/docs/cli/get-started/#step-2-turn-on-the-1password-desktop-app-integration\n\n")
				return "", false
			}

			if out, err := exec.CommandContextQ(ctx, "op", "read", ref); err == nil {
				return strings.TrimSpace(string(out)), true
			}

			return "", false
		}, nil
	},
}
