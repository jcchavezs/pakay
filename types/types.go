package types

import "context"

type (
	SecretGetter   func(ctx context.Context) (string, bool)
	SecretProvider func(cfg ProviderConfig) (SecretGetter, error)
	ProviderConfig map[string]any
)
