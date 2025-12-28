package pakay

import (
	internaltypes "github.com/jcchavezs/pakay/internal/types"
	"github.com/jcchavezs/pakay/types"

	"github.com/jcchavezs/pakay/internal/parser"
)

// SecretSource describes how to retrieve a secret from a given backend and the labels
// that should be applied to the secret entry.
type SecretSource struct {
	internaltypes.TypedConfig
	Labels []string
}

// SecretConfig represents a single secret definition including its sources and metadata.
type SecretConfig struct {
	Name        string
	Description string
	Sources     []SecretSource
}

// SecretsConfig groups multiple secret definitions that together form a manifest.
type SecretsConfig []SecretConfig

// toManifestEntries converts the in-memory secrets config into parser manifest entries.
func (ssc SecretsConfig) toManifestEntries() []parser.ManifestEntry {
	entries := make([]parser.ManifestEntry, 0, len(ssc))
	for _, sc := range ssc {
		me := parser.ManifestEntry{
			Name:        sc.Name,
			Description: sc.Description,
			Sources:     make([]parser.ManifestEntrySource, 0, len(sc.Sources)),
		}

		for _, s := range sc.Sources {
			c := s.TypedConfig.(types.SourceConfig)
			me.Sources = append(me.Sources, parser.ManifestEntrySource{
				Labels: s.Labels,
				Type:   c.Type(),
				Config: c,
			})
		}

		entries = append(entries, me)
	}
	return entries
}
