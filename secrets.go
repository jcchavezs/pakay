package pakay

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/jcchavezs/pakay/internal/providers/env"
	onepasswordcli "github.com/jcchavezs/pakay/internal/providers/onepassword/cli"
	"github.com/jcchavezs/pakay/internal/providers/stdin"
	"github.com/jcchavezs/pakay/internal/values"
	"github.com/jcchavezs/pakay/types"
)

func init() {
	RegisterProvider("stdin", stdin.Provider)
	RegisterProvider("env", env.Provider)
	RegisterProvider("onepassword_cli", onepasswordcli.Provider)
}

var providers = map[string]types.SecretProvider{}

// RegisterProvider registers a new secret provider with the given name.
// If a provider with the same name already exists, it will be overwritten.
// This function is intended to be used by package authors to add their own secret providers.
func RegisterProvider(name string, p types.SecretProvider) {
	providers[name] = p
}

type (
	secret struct {
		manifestEntry
		getters []types.SecretGetter
	}
)

var secrets = map[string]secret{}

type LoadOptions struct {
	Variables map[string]string
}

// LoadSecretsFromBytes loads secrets from a YAML manifest provided as a byte slice.
// The manifest should contain a list of secrets with their names, descriptions, and sources.
// Each source should specify a type and its configuration.
func LoadSecretsFromBytes(config []byte, opts LoadOptions) error {
	cfg, err := parseManifest(config, opts)
	if err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	for _, c := range cfg {
		s := secret{
			manifestEntry: c,
			getters:       make([]types.SecretGetter, 0, len(c.Sources)),
		}

		for _, src := range c.Sources {
			t, ok := values.GetFromMap[string](src, "type")
			if !ok {
				continue
			}

			p, ok := providers[t]
			if !ok {
				return fmt.Errorf("unknown provider: %s", t)
			}

			tc, ok := values.GetFromMap[map[string]any](src, t)
			if !ok {
				return fmt.Errorf("missing configuration for provider: %s", t)
			}

			g, err := p(tc)
			if err != nil {
				return fmt.Errorf("building secret getter: %w", err)
			}

			s.getters = append(s.getters, g)
		}

		secrets[c.Name] = s
	}

	return nil
}

// GetSecret retrieves the value of a secret by its name.
// It returns the secret value and a boolean indicating whether the secret was found.
// If the secret is not found, it logs an error and returns an empty string and false.
// The function will try each getter associated with the secret until it finds a valid value.
// If no getter returns a valid value, it will return an empty string and false.
func GetSecret(ctx context.Context, name string) (string, bool) {
	s, ok := secrets[name]
	if !ok {
		slog.Error("Secret not found.", "name", name)
		return "", false
	}

	for _, g := range s.getters {
		if val, ok := g(ctx); ok {
			return val, true
		}
	}

	return "", false
}

type manifestEntry struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Sources     []map[string]any `yaml:"sources"`
}

// parseManifest parses the YAML manifest and returns a slice of manifestEntry.
// If the manifest contains variables, it will render them using the provided LoadOptions.
// It returns an error if the manifest cannot be parsed or rendered.
func parseManifest(manifest []byte, opt LoadOptions) ([]manifestEntry, error) {
	var cfg []manifestEntry

	rConfig := manifest
	if len(opt.Variables) > 0 {
		tmpl, err := template.New("manifest").Parse(string(manifest))
		if err != nil {
			return nil, fmt.Errorf("parsing manifest: %w", err)
		}

		s := bytes.Buffer{}
		if err = tmpl.Execute(&s, opt.Variables); err != nil {
			return nil, fmt.Errorf("rendering manifest: %w", err)
		}

		rConfig = s.Bytes()
	}

	if err := yaml.Unmarshal(rConfig, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling manifest: %w", err)
	}

	return cfg, nil
}
