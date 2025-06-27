package static

import (
	"context"
	"errors"

	"github.com/jcchavezs/pakay/types"
)

type Config struct {
	Value string `yaml:"value"`
}

func (c Config) String() string {
	return c.Value
}

var Provider = types.SecretProvider{
	ConfigFactory: func() types.ProviderConfig {
		return &Config{}
	},
	SecretGetterFactory: func(cfg types.ProviderConfig) (types.SecretGetter, error) {
		var val string
		if tCfg, ok := cfg.(*Config); ok {
			val = tCfg.Value
		} else {
			return nil, errors.New("invalid config")
		}

		return func(context.Context) (string, bool) {
			return val, true
		}, nil
	},
}
