package secrets

import (
	"github.com/jcchavezs/pakay/internal/parser"
	"github.com/jcchavezs/pakay/types"
)

type (
	Getter struct {
		Labels []string
		types.SecretGetter
	}

	Secret struct {
		parser.ManifestEntry
		Getters []Getter
	}
)

var (
	All    = map[string]Secret{}
	Loaded bool
)
