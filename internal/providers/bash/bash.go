package bash

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"time"

	"github.com/jcchavezs/pakay/types"
)

type Config struct {
	Command   string `yaml:"command"`
	TimeoutMS int    `yaml:"timeout_ms"`
}

func (c Config) String() string {
	return c.Command
}

var Provider = types.SecretProvider{
	ConfigFactory: func() types.ProviderConfig {
		return &Config{}
	},
	SecretGetterFactory: func(cfg types.ProviderConfig) (types.SecretGetter, error) {
		var (
			command string
			timeout time.Duration
		)
		if tCfg, ok := cfg.(*Config); ok {
			command = tCfg.Command
			timeout = time.Duration(tCfg.TimeoutMS) * time.Millisecond
		} else {
			return nil, errors.New("invalid config")
		}

		return func(ctx context.Context) (string, bool) {
			if timeout > 0 {
				var cancelFn context.CancelFunc
				ctx, cancelFn = context.WithTimeout(ctx, timeout)
				defer cancelFn()
			}

			cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
			cmd.Stderr = os.Stderr
			out, err := cmd.Output()
			if err != nil {
				return "", false
			}

			return string(bytes.TrimSpace(out)), true
		}, nil
	},
}
