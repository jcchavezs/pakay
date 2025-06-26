package secrets

import (
	"github.com/jcchavezs/pakay/internal/parser"
	"github.com/jcchavezs/pakay/types"
)

type (
	Secret struct {
		parser.ManifestEntry
		Getters []types.SecretGetter
	}
)

var (
	All    = map[string]Secret{}
	Loaded bool
)
