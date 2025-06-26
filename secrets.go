package pakay

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jcchavezs/pakay/internal/log"
	"github.com/jcchavezs/pakay/internal/parser"
	"github.com/jcchavezs/pakay/internal/providers"
	"github.com/jcchavezs/pakay/internal/secrets"
	"github.com/jcchavezs/pakay/types"
)

// RegisterProvider registers a new secret provider with the given name.
// If a provider with the same name already exists, it will be overwritten.
// This function is intended to be used by package authors to add their own secret providers.
var RegisterProvider = providers.RegisterProvider

type LoadOptions struct {
	Variables  map[string]string
	LogHandler slog.Handler
}

// LoadSecretsFromBytes loads secrets from a YAML manifest provided as a byte slice.
// The manifest should contain a list of secrets with their names, descriptions, and sources.
// Each source should specify a type and its configuration.
func LoadSecretsFromBytes(config []byte) error {
	return LoadSecretsFromBytesWithOptions(config, LoadOptions{})
}

func LoadSecretsFromBytesWithOptions(config []byte, opts LoadOptions) error {
	cfg, err := parser.ParseManifest(config, opts.Variables)
	if err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	for _, c := range cfg {
		s := secrets.Secret{
			ManifestEntry: c,
			Getters:       make([]types.SecretGetter, 0, len(c.Sources)),
		}

		for _, src := range c.Sources {
			p, ok := providers.GetProvider(src.Type)
			if !ok {
				return fmt.Errorf("unknown provider: %s", src.Type)
			}

			g, err := p.SecretGetterFactory(src.Config)
			if err != nil {
				return fmt.Errorf("building secret getter: %w", err)
			}

			s.Getters = append(s.Getters, g)
		}

		secrets.All[c.Name] = s
	}

	if opts.LogHandler != nil {
		log.SetHandler(opts.LogHandler)
	}

	secrets.Loaded = true

	return nil
}

// GetSecret retrieves the value of a secret by its name.
// It returns the secret value and a boolean indicating whether the secret was found.
// If the secret is not found, it logs an error and returns an empty string and false.
// The function will try each getter associated with the secret until it finds a valid value.
// If no getter returns a valid value, it will return an empty string and false.
func GetSecret(ctx context.Context, name string) (string, bool) {
	s, ok := secrets.All[name]
	if !ok {
		return "", false
	}

	for _, g := range s.Getters {
		if val, ok := g(ctx); ok {
			return val, true
		}
	}

	return "", false
}

// AssertSecrets asserts the availability of the loaded secrets.
// It is useful to check the secrets before running the command.
func AssertSecrets(ctx context.Context) ([]string, error) {
	if !secrets.Loaded {
		return nil, errors.New("secrets haven't been loaded")
	}

	missing := []string{}
	for name := range secrets.All {
		if _, ok := GetSecret(ctx, name); !ok {
			missing = append(missing, name)
		}
	}

	return missing, nil
}
