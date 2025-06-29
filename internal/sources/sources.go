package sources

import (
	"github.com/jcchavezs/pakay/internal/sources/bash"
	"github.com/jcchavezs/pakay/internal/sources/env"
	onepasswordcli "github.com/jcchavezs/pakay/internal/sources/onepassword/cli"
	"github.com/jcchavezs/pakay/internal/sources/static"
	"github.com/jcchavezs/pakay/internal/sources/stdin"
	"github.com/jcchavezs/pakay/types"
)

var sources = map[string]types.SecretSource{}

func Register(name string, p types.SecretSource) {
	sources[name] = p
}

func Get(name string) (types.SecretSource, bool) {
	p, ok := sources[name]
	return p, ok
}

func init() {
	Register("static", static.Source)
	Register("bash", bash.Source)
	Register("env", env.Source)
	Register("stdin", stdin.Source)
	Register("1password", onepasswordcli.Source)
}
