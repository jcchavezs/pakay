package types

import (
	"context"
	"fmt"
)

type (
	SecretGetter   func(ctx context.Context) (string, bool)
	SecretProvider struct {
		ConfigFactory       func() ProviderConfig
		SecretGetterFactory func(cfg ProviderConfig) (SecretGetter, error)
	}
	ProviderConfig fmt.Stringer
)
