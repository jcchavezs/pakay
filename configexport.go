//go:generate go run ./internal/cmd/exportconfig

package pakay

import (
	"github.com/jcchavezs/pakay/internal/sources/onepassword/cli"
	"github.com/jcchavezs/pakay/internal/sources/static"
	"github.com/jcchavezs/pakay/internal/sources/bash"
	"github.com/jcchavezs/pakay/internal/sources/env"
	"github.com/jcchavezs/pakay/internal/sources/stdin"
)

type (
	StaticConfig = static.Config
	BashConfig = bash.Config
	EnvConfig = env.Config
	StdinConfig = stdin.Config
	OnePasswordConfig = cli.Config
)
