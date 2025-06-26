package env

import (
	"context"
	"errors"
	"os"

	"github.com/jcchavezs/pakay/types"
)

type Config struct {
	Name string `yaml:"name"`
}

func (c Config) String() string {
	return c.Name
}

var Provider = types.SecretProvider{
	ConfigFactory: func() types.ProviderConfig {
		return &Config{}
	},
	SecretGetterFactory: func(cfg types.ProviderConfig) (types.SecretGetter, error) {
		var name string
		if tCfg, ok := cfg.(*Config); ok {
			name = tCfg.Name
		} else {
			return nil, errors.New("invalid config")
		}

		return func(context.Context) (string, bool) {
			val := os.Getenv(name)
			return val, val != ""
		}, nil
	},
}
