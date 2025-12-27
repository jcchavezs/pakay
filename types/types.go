package types

import (
	"context"
	"fmt"

	internaltypes "github.com/jcchavezs/pakay/internal/types"
)

type (
	// SecretGetter gets a given secret
	SecretGetter func(ctx context.Context) (string, bool)

	// SecretSource is a source for a given secret
	SecretSource struct {
		ConfigFactory       func() SourceConfig
		SecretGetterFactory func(cfg SourceConfig) (SecretGetter, error)
	}

	// SourceConfig is the config for a source of a given secret
	SourceConfig interface {
		fmt.Stringer
		internaltypes.TypedConfig
		Type() string
	}
)
