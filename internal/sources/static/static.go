package static

import (
	"context"
	"errors"
	"strings"

	internaltypes "github.com/jcchavezs/pakay/internal/types"

	"github.com/jcchavezs/pakay/types"
)

type Config struct {
	Value string `yaml:"value"`
}

func (c *Config) String() string {
	l := len(c.Value)
	if l == 0 {
		return ""
	}

	hidden := l
	if l > 4 {
		hidden = l - 3
	}

	return string(c.Value[0:l-hidden]) + strings.Repeat("*", hidden)
}

func (*Config) Type() string {
	return "static"
}

func (c *Config) SentinelFn() internaltypes.SentinelVal { return internaltypes.SentinelVal{} }

var Source = types.SecretSource{
	ConfigFactory: func() types.SourceConfig {
		return &Config{}
	},
	SecretGetterFactory: func(cfg types.SourceConfig) (types.SecretGetter, error) {
		var val string
		if tCfg, ok := cfg.(*Config); ok {
			val = tCfg.Value
		} else {
			return nil, errors.New("invalid config")
		}

		if val == "" {
			return nil, errors.New("value cannot be empty")
		}

		return func(context.Context) (string, bool) {
			return val, true
		}, nil
	},
}
