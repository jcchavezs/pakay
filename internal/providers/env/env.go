package env

import (
	"context"
	"errors"
	"os"

	"github.com/jcchavezs/pakay/internal/values"
	"github.com/jcchavezs/pakay/types"
)

var Provider types.SecretProvider = func(cfg types.ProviderConfig) (types.SecretGetter, error) {
	var name, found = values.GetFromMap[string](cfg, "name")
	if !found {
		return nil, errors.New("missing env.name value")
	}

	return func(context.Context) (string, bool) {
		val := os.Getenv(name)
		return val, val != ""
	}, nil
}
