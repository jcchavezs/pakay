package env

import (
	"context"
	"errors"
	"os"

	internaltypes "github.com/jcchavezs/pakay/internal/types"

	"github.com/jcchavezs/pakay/types"
)

type Config struct {
	internaltypes.TypedConfig
	Key string `yaml:"key"`
}

func (c *Config) String() string {
	return c.Key
}

func (*Config) Type() string {
	return "env"
}

func (*Config) SentinelFn(internaltypes.SentinelVal) {}

var Source = types.SecretSource{
	ConfigFactory: func() types.SourceConfig {
		return &Config{}
	},
	SecretGetterFactory: func(cfg types.SourceConfig) (types.SecretGetter, error) {
		var key string
		if tCfg, ok := cfg.(*Config); ok {
			key = tCfg.Key
		} else {
			return nil, errors.New("invalid config")
		}

		if key == "" {
			return nil, errors.New("key cannot be empty")
		}

		return func(context.Context) (string, bool) {
			val := os.Getenv(key)
			return val, val != ""
		}, nil
	},
}
