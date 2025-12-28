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

// ParseAndLoadOptions defines options for loading secrets.
type ParseAndLoadOptions struct {
	Variables map[string]string
	LoadOptions
}

// LoadOptions defines options for loading secrets.
type LoadOptions struct {
	LogHandler slog.Handler
}

// LoadSecrets loads secrets from the provided SecretsConfig.
func LoadSecrets(config SecretsConfig) error {
	return loadSecretsFromManifestEntries(config.toManifestEntries(), LoadOptions{})
}

// LoadSecretsWithOptions loads secrets from the provided SecretsConfig with additional options.
func LoadSecretsWithOptions(config SecretsConfig, opts LoadOptions) error {
	return loadSecretsFromManifestEntries(config.toManifestEntries(), opts)
}

var sMutex sync.RWMutex

func loadSecretsFromManifestEntries(cfg []parser.ManifestEntry, opts LoadOptions) error {
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
				return fmt.Errorf("building secret getter for %s: %w", p.ConfigFactory().Type(), err)
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
		log.SetHandler(log.DiscardHandler)
	} else {
		log.SetHandler(opts.LogHandler)
	}

	sMutex.Lock()
	secrets.Loaded = true
	sMutex.Unlock()

	return nil
}

// ParseAndLoadSecrets loads secrets from a YAML manifest provided as a byte slice.
// The manifest should contain a list of secrets with their names, descriptions, and sources.
// Each source should specify a type and its configuration.
func ParseAndLoadSecrets(manifest []byte) error {
	return ParseAndLoadSecretsWithOptions(manifest, ParseAndLoadOptions{})
}

// ParseAndLoadSecretsWithOptions loads secrets from a YAML manifest provided as a byte slice with additional options.
// The manifest should contain a list of secrets with their names, descriptions, and sources.
// Each source should specify a type and its configuration.
func ParseAndLoadSecretsWithOptions(manifest []byte, opts ParseAndLoadOptions) error {
	cfg, err := parser.ParseManifest(manifest, opts.Variables)
	if err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	return loadSecretsFromManifestEntries(cfg, opts.LoadOptions)
}

// GetSecret retrieves the value of a secret by its name.
// It returns the secret value and a boolean indicating whether the secret was found.
// If the secret is not found, it logs an error and returns an empty string and false.
// The function will try each getter associated with the secret until it finds a valid value.
// If no getter returns a valid value, it will return an empty string and false.
func GetSecret(ctx context.Context, name string) (string, bool) {
	return GetSecretWithOptions(ctx, name, SecretOptions{})
}

// SecretOptions defines options for retrieving secrets.
type SecretOptions struct {
	FilterIn FilterIn[Source]
}

// GetSecretWithOptions retrieves the value of a secret by its name with additional options.
// It returns the secret value and a boolean indicating whether the secret was found.
// If the secret is not found, it logs an error and returns an empty string and false.
// The function will try each getter associated with the secret that passes the FilterIn function
// until it finds a valid value. If no getter returns a valid value, it will return an empty string and false.
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

func checkSecretsAreLoaded() bool {
	sMutex.RLock()
	defer sMutex.RUnlock()
	return secrets.Loaded
}

// AssertSecrets asserts the availability of the loaded secrets.
// It is useful to check the secrets before running the command.
func AssertSecrets(ctx context.Context) ([]string, error) {
	return AssertSecretsWithOptions(ctx, AssertOptions{})
}

// AssertOptions defines options for asserting secrets.
type AssertOptions struct {
	SecretFilterIn FilterIn[Secret]
	SourceFilterIn FilterIn[Source]
}

// Secret represents a secret with its name.
type Secret struct {
	Name string
}

// AssertSecrets asserts the availability of the loaded secrets.
// It is useful to check the secrets before running the command.
func AssertSecretsWithOptions(ctx context.Context, opts AssertOptions) ([]string, error) {
	if !checkSecretsAreLoaded() {
		return nil, errors.New("secrets haven't been loaded yet")
	}

	missing := []string{}
	for name := range secrets.All {
		if opts.SecretFilterIn != nil {
			if !opts.SecretFilterIn(Secret{Name: name}) {
				continue
			}
		}

		if _, ok := GetSecretWithOptions(ctx, name, SecretOptions{FilterIn: opts.SourceFilterIn}); !ok {
			missing = append(missing, name)
		}
	}

	return missing, nil
}
