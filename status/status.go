package status

import (
	"context"

	"github.com/jcchavezs/pakay"
	"github.com/jcchavezs/pakay/internal/secrets"
)

type SecretStatus struct {
	secret    secrets.Secret
	available bool
}

func (ss SecretStatus) Name() string {
	return ss.secret.Name
}

func (ss SecretStatus) Description() string {
	return ss.secret.Description
}

func (ss SecretStatus) Sources() []string {
	sources := make([]string, 0, len(ss.secret.Sources))
	for _, s := range ss.secret.Sources {
		sources = append(sources, s.String())
	}

	return sources
}

func (ss SecretStatus) Available() bool {
	return ss.available
}

func CheckSecrets(ctx context.Context) ([]SecretStatus, bool) {
	ss := make([]SecretStatus, 0, len(secrets.All))
	allAvailable := true

	for name, s := range secrets.All {
		_, ok := pakay.GetSecret(ctx, name)
		allAvailable = allAvailable && ok
		ss = append(ss, SecretStatus{
			secret:    s,
			available: ok,
		})
	}

	return ss, allAvailable
}
