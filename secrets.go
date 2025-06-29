package pakay

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jcchavezs/pakay/internal/log"
	"github.com/jcchavezs/pakay/internal/parser"
	"github.com/jcchavezs/pakay/internal/secrets"
	"github.com/jcchavezs/pakay/internal/sources"
)

// RegisterSource registers a new secret source with the given name.
// If a source with the same name already exists, it will be overwritten.
// This function is intended to be used by package authors to add their own secret sources.
var RegisterSource = sources.Register

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

var sMutex sync.RWMutex

func LoadSecretsFromBytesWithOptions(config []byte, opts LoadOptions) error {
	cfg, err := parser.ParseManifest(config, opts.Variables)
	if err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	for _, c := range cfg {
		if _, ok := secrets.All[c.Name]; ok {
			return fmt.Errorf("duplicated declaration for %q", c.Name)
		}

		s := secrets.Secret{
			ManifestEntry: c,
			Getters:       make([]secrets.Getter, 0, len(c.Sources)),
		}

		for _, src := range c.Sources {
			p, ok := sources.Get(src.Type)
			if !ok {
				return fmt.Errorf("unknown source: %s", src.Type)
			}

			g, err := p.SecretGetterFactory(src.Config)
			if err != nil {
				return fmt.Errorf("building secret getter: %w", err)
			}

			sMutex.Lock()
			s.Getters = append(s.Getters, secrets.Getter{
				Labels:       src.Labels,
				SecretGetter: g,
			})
			sMutex.Unlock()
		}

		secrets.All[c.Name] = s
	}

	if opts.LogHandler == nil {
		log.SetHandler(DiscardHandler)
	} else {
		log.SetHandler(opts.LogHandler)
	}

	sMutex.Lock()
	secrets.Loaded = true
	sMutex.Unlock()

	return nil
}

// GetSecret retrieves the value of a secret by its name.
// It returns the secret value and a boolean indicating whether the secret was found.
// If the secret is not found, it logs an error and returns an empty string and false.
// The function will try each getter associated with the secret until it finds a valid value.
// If no getter returns a valid value, it will return an empty string and false.
func GetSecret(ctx context.Context, name string) (string, bool) {
	return GetSecretWithOptions(ctx, name, SecretOptions{})
}

type SecretOptions struct {
	FilterIn FilterIn
}

func checkSecretsAreLoaded() bool {
	sMutex.RLock()
	defer sMutex.RUnlock()
	return secrets.Loaded
}

func GetSecretWithOptions(ctx context.Context, name string, opts SecretOptions) (string, bool) {
	if !checkSecretsAreLoaded() {
		log.Logger.Error("Secrets haven't been loaded yet")
		return "", false
	}

	s, ok := secrets.All[name]
	if !ok {
		log.Logger.Error("Unknown secret", "name", name)
		return "", false
	}

	for i, g := range s.Getters {
		if opts.FilterIn != nil {
			if !opts.FilterIn(Source{Type: s.ManifestEntry.Sources[i].Type, Labels: g.Labels}) {
				continue
			}
		}

		if val, ok := g.SecretGetter(ctx); ok {
			return val, true
		}
	}

	return "", false
}

func AssertSecrets(ctx context.Context) ([]string, error) {
	return AssertSecretsWithOptions(ctx, AssertOptions{})
}

type AssertOptions struct {
	FilterIn FilterIn
}

// AssertSecrets asserts the availability of the loaded secrets.
// It is useful to check the secrets before running the command.
func AssertSecretsWithOptions(ctx context.Context, opts AssertOptions) ([]string, error) {
	if !checkSecretsAreLoaded() {
		return nil, errors.New("secrets haven't been loaded yet")
	}

	missing := []string{}
	for name := range secrets.All {
		if _, ok := GetSecretWithOptions(ctx, name, (SecretOptions)(opts)); !ok {
			missing = append(missing, name)
		}
	}

	return missing, nil
}
