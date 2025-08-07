package sources

import (
	"maps"
	"slices"

	"github.com/jcchavezs/pakay/internal/sources/bash"
	"github.com/jcchavezs/pakay/internal/sources/env"
	onepasswordcli "github.com/jcchavezs/pakay/internal/sources/onepassword/cli"
	"github.com/jcchavezs/pakay/internal/sources/static"
	"github.com/jcchavezs/pakay/internal/sources/stdin"
	"github.com/jcchavezs/pakay/types"
)

var (
	sources = map[string]types.SecretSource{}
)

func Register(p types.SecretSource) {
	sources[p.ConfigFactory().Type()] = p
}

func Get(name string) (types.SecretSource, bool) {
	p, ok := sources[name]
	return p, ok
}

func GetAll() []types.SecretSource {
	return slices.Collect(maps.Values(sources))
}

func init() {
	Register(static.Source)
	Register(bash.Source)
	Register(env.Source)
	Register(stdin.Source)
	Register(onepasswordcli.Source)
}
