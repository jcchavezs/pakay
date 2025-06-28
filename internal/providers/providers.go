package providers

import (
	"github.com/jcchavezs/pakay/internal/providers/bash"
	"github.com/jcchavezs/pakay/internal/providers/env"
	onepasswordcli "github.com/jcchavezs/pakay/internal/providers/onepassword/cli"
	"github.com/jcchavezs/pakay/internal/providers/static"
	"github.com/jcchavezs/pakay/internal/providers/stdin"
	"github.com/jcchavezs/pakay/types"
)

var providers = map[string]types.SecretProvider{}

func RegisterProvider(name string, p types.SecretProvider) {
	providers[name] = p
}

func GetProvider(name string) (types.SecretProvider, bool) {
	p, ok := providers[name]
	return p, ok
}

func init() {
	RegisterProvider("static", static.Provider)
	RegisterProvider("bash", bash.Provider)
	RegisterProvider("env", env.Provider)
	RegisterProvider("stdin", stdin.Provider)
	RegisterProvider("1password", onepasswordcli.Provider)
}
