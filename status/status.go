package status

import (
	"context"

	"github.com/jcchavezs/pakay"
	"github.com/jcchavezs/pakay/internal/secrets"
)

type SecretStatus struct {
	filterIn  pakay.FilterIn
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
		if ss.filterIn != nil {
			if !ss.filterIn(pakay.Source{Type: s.Type, Labels: s.Labels}) {
				continue
			}
		}

		sources = append(sources, s.String())
	}

	return sources
}

func (ss SecretStatus) Available() bool {
	return ss.available
}

type CheckOptions struct {
	FilterIn pakay.FilterIn
}

func CheckSecrets(ctx context.Context) []SecretStatus {
	return CheckSecretsWithOptions(ctx, CheckOptions{})
}

func CheckSecretsWithOptions(ctx context.Context, opts CheckOptions) []SecretStatus {
	ss := make([]SecretStatus, 0, len(secrets.All))

	for name, s := range secrets.All {
		_, ok := pakay.GetSecretWithOptions(ctx, name, pakay.SecretOptions{
			FilterIn: opts.FilterIn,
		})
		ss = append(ss, SecretStatus{
			filterIn:  opts.FilterIn,
			secret:    s,
			available: ok,
		})
	}

	return ss
}
