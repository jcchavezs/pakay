package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/goccy/go-yaml"
	"github.com/jcchavezs/pakay/internal/providers"
	"github.com/jcchavezs/pakay/types"
)

type ManifestEntrySource struct {
	Type   string `yaml:"type"`
	Config types.ProviderConfig
}

func (s ManifestEntrySource) String() string {
	return fmt.Sprintf("%s: %s", s.Type, s.Config)
}

func (s *ManifestEntrySource) UnmarshalYAML(data []byte) error {
	t := struct {
		Type string `yaml:"type"`
	}{}
	if err := yaml.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("unmarshaling type: %w", err)
	}

	p, ok := providers.GetProvider(t.Type)
	if !ok {
		return fmt.Errorf("unknown provider: %s", t.Type)
	}

	cfg := map[string]json.RawMessage{}
	if err := yaml.UnmarshalWithOptions(data, &cfg, yaml.UseJSONUnmarshaler()); err != nil {
		return fmt.Errorf("unmarshaling provider raw configuration: %w", err)
	}

	pcfg, ok := cfg[t.Type]
	if !ok {
		return fmt.Errorf("missing provider configuration")
	}

	tCfg := p.ConfigFactory()
	if err := yaml.Unmarshal([]byte(pcfg), tCfg); err != nil {
		return fmt.Errorf("unmarshaling provider typed configuration: %w", err)
	}

	s.Type = t.Type
	s.Config = tCfg

	return nil
}

type ManifestEntry struct {
	Name        string                `yaml:"name"`
	Description string                `yaml:"description"`
	Sources     []ManifestEntrySource `yaml:"sources"`
}

// parseManifest parses the YAML manifest and returns a slice of manifestEntry.
// If the manifest contains variables, it will render them using the provided LoadOptions.
// It returns an error if the manifest cannot be parsed or rendered.
func ParseManifest(manifest []byte, vars map[string]string) ([]ManifestEntry, error) {
	var cfg []ManifestEntry

	rConfig := manifest
	if len(vars) > 0 {
		tmpl, err := template.New("manifest").Parse(string(manifest))
		if err != nil {
			return nil, fmt.Errorf("parsing manifest: %w", err)
		}

		s := bytes.Buffer{}
		if err = tmpl.Execute(&s, vars); err != nil {
			return nil, fmt.Errorf("rendering manifest: %w", err)
		}

		rConfig = s.Bytes()
	}

	if err := yaml.UnmarshalWithOptions(rConfig, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling manifest: %w", err)
	}

	return cfg, nil
}
